package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/futugyou/yomawari/mcp"
	"github.com/futugyou/yomawari/mcp/protocol/messages"
	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/runtime/tasks"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var s_clientSessionDuration = mcp.CreateDurationHistogram("mcp.client.session.duration", "Measures the duration of a client session.", true)
var s_serverSessionDuration = mcp.CreateDurationHistogram("mcp.server.session.duration", "Measures the duration of a server session.", true)
var s_clientOperationDuration = mcp.CreateDurationHistogram("mcp.client.operation.duration", "Measures the duration of outbound message.", false)
var s_serverOperationDuration = mcp.CreateDurationHistogram("mcp.server.operation.duration", "Measures the duration of inbound message processing.", false)

type McpSession struct {
	_isServer                 bool
	_transportKind            string
	_transport                transport.ITransport
	_requestHandlers          *RequestHandlers
	_notificationHandlers     *NotificationHandlers
	_sessionStartingTimestamp int64
	_pendingRequests          sync.Map // map[RequestId]*tasks.TaskCompletionSource[messages.IJsonRpcMessage]
	_handlingRequests         sync.Map // map[RequestId]context.CancelFunc
	_id                       string
	_nextRequestId            int64
	EndpointName              string
}

func NewMcpSession(isServer bool, transp transport.ITransport, endpointName string, requestHandlers *RequestHandlers, notificationHandlers *NotificationHandlers) *McpSession {
	_transportKind := "unknownTransport"
	switch transp.(type) {
	case *transport.StdioClientSessionTransport, *transport.StdioServerTransport:
		_transportKind = "stdio"
	case *transport.StreamClientSessionTransport, *transport.StreamServerTransport:
		_transportKind = "stream"
	case *transport.SseClientSessionTransport, *transport.SseResponseStreamTransport:
		_transportKind = "sse"
	}

	return &McpSession{
		_isServer:                 isServer,
		_transportKind:            _transportKind,
		_transport:                transp,
		_requestHandlers:          requestHandlers,
		_notificationHandlers:     notificationHandlers,
		_sessionStartingTimestamp: time.Now().UnixNano(),
		_id:                       uuid.New().String(),
		_nextRequestId:            0,
		EndpointName:              endpointName,
	}
}

func (m *McpSession) ProcessMessages(ctx context.Context) error {
	var processMessage = func(ctx context.Context, message messages.IJsonRpcMessage) error {
		var messageWithId messages.IJsonRpcMessageWithId
		if msg, ok := message.(messages.IJsonRpcMessageWithId); ok {
			messageWithId = msg
		}

		var combinedCts context.Context
		var cancelfunc context.CancelFunc

		if messageWithId != nil && messageWithId.GetId() != nil {
			combinedCts, cancelfunc = context.WithCancel(ctx)
			id := messageWithId.GetId()
			m._handlingRequests.Store(*id, cancelfunc)
		}

		// maybe need maybe not
		runtime.Gosched()

		callCtx := ctx
		if combinedCts != nil {
			callCtx = combinedCts
		}

		if err := m.handleMessage(callCtx, message); err != nil {
			isUserCancellation := false
			if (err == context.Canceled || err == context.DeadlineExceeded) && callCtx.Err() != nil {
				isUserCancellation = true
			}
			if request, ok := message.(*messages.JsonRpcRequest); !isUserCancellation && ok {
				m._transport.SendMessage(ctx, messages.NewJsonRpcError(request.Id, 500, err.Error(), nil))
			}
		}
		return nil
	}

	resultChan := m._transport.MessageReader()
	processor := tasks.TaskProcessor[messages.IJsonRpcMessage]{
		ResultChan:        resultChan,
		Handler:           processMessage, // func(ctx, msg) error
		MaxConcurrency:    20,
		PerMessageTimeout: 10 * time.Second,
	}

	processor.Run(ctx)

	m._pendingRequests.Range(func(key, value interface{}) bool {
		if tcs, ok := value.(*tasks.TaskCompletionSource[messages.IJsonRpcMessage]); ok {
			tcs.TrySetError(fmt.Errorf("the server shut down unexpectedly"))
		}
		return true
	})

	return nil
}

func (m *McpSession) handleMessage(ctx context.Context, message messages.IJsonRpcMessage) error {
	durationMetric := s_clientOperationDuration
	if m._isServer {
		durationMetric = s_serverOperationDuration
	}
	method := getMethodName(message)

	var startingTimestamp int64 = time.Now().UnixNano()
	ctx, span := mcp.StartSpanWithJsonRpcData(ctx, m.createActivityName(method), message)
	defer span.End()

	tags := []attribute.KeyValue{}
	m.addStandardTags(&tags, method)
	defer finalizeDiagnostics(ctx, &startingTimestamp, durationMetric, tags)

	if err := m._transport.SendMessage(ctx, message); err != nil {
		addExceptionTags(&tags, err)
		return err
	}

	var err error
	switch request := message.(type) {
	case *messages.JsonRpcRequest:
		addRpcRequestTags(&tags, *request)
		err = m.handleRequest(ctx, *request)
	case *messages.JsonRpcNotification:
		err = m.handleNotification(ctx, request)
	case messages.IJsonRpcMessageWithId:
		err = m.handleMessageWithId(message, request)
	default:
	}

	if err != nil {
		addExceptionTags(&tags, err)
		return err
	}

	return nil
}

func (m *McpSession) handleNotification(ctx context.Context, notification *messages.JsonRpcNotification) error {
	if notification.Method == messages.NotificationMethods_CancelledNotification {
		if cn := getCancelledNotificationParams(notification.Params); cn != nil {
			value, ok := m._handlingRequests.Load(cn.RequestId)
			if ok {
				if cancel, ok := value.(context.CancelFunc); ok {
					cancel()
				}
			}
		}
	}

	return m._notificationHandlers.InvokeHandlers(ctx, notification.Method, notification)
}

func (m *McpSession) handleMessageWithId(message messages.IJsonRpcMessage, messageWithId messages.IJsonRpcMessageWithId) error {
	if messageWithId.GetId() == nil || len(messageWithId.GetId().String()) == 0 {
		return fmt.Errorf("message with id has no id")
	}
	requestid := messageWithId.GetId()
	value, ok := m._pendingRequests.Load(*requestid)
	if !ok {
		return fmt.Errorf("no pending request found for id %s", requestid.String())
	}
	if source, ok := value.(*tasks.TaskCompletionSource[messages.IJsonRpcMessage]); ok {
		source.SetResult(message)
	}
	return nil
}

func (m *McpSession) handleRequest(ctx context.Context, request messages.JsonRpcRequest) error {
	handler, ok := m._requestHandlers.Get(request.Method)
	if !ok {
		return fmt.Errorf("no handler found for method %s", request.Method)
	}

	result, err := handler(ctx, &request)
	if err != nil {
		return err
	}

	msg := messages.NewJsonRpcResponse(request.Id, result)
	return m._transport.SendMessage(ctx, msg)
}

func (m *McpSession) RegisterNotificationHandler(method string, handler NotificationHandler) *RegistrationHandle {
	return m._notificationHandlers.Register(method, handler, true)
}

func (m *McpSession) createActivityName(method string) string {
	s := "client"
	if m._isServer {
		s = "server"
	}

	return fmt.Sprintf("mcp.%s.%s/%s", s, m._transportKind, method)
}

func (m *McpSession) SendRequest(ctx context.Context, request *messages.JsonRpcRequest) (*messages.JsonRpcResponse, error) {
	if m == nil || request == nil {
		return nil, fmt.Errorf("session or request is nil")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	durationMetric := s_clientOperationDuration
	if m._isServer {
		durationMetric = s_serverOperationDuration
	}

	method := request.Method

	var startingTimestamp int64 = time.Now().UnixNano()
	ctx, span := mcp.Tracer.Start(ctx, m.createActivityName(method))
	defer span.End()

	tags := []attribute.KeyValue{}
	m.addStandardTags(&tags, method)
	defer finalizeDiagnostics(ctx, &startingTimestamp, durationMetric, tags)

	// Set request ID
	if request.Id == nil {
		newId := atomic.AddInt64(&m._nextRequestId, 1)
		request.Id = messages.NewRequestIdFromString(fmt.Sprintf("%s-%d", m._id, newId))
	}

	mcp.PropagatorInject(ctx, request)

	ctx, cancelfunc := context.WithCancel(ctx)
	tcs := tasks.NewTaskCompletionSource[messages.IJsonRpcMessage](ctx, cancelfunc)
	m._pendingRequests.Store(request.Id, tcs)

	m.addStandardTags(&tags, method)
	addRpcRequestTags(&tags, *request)

	defer finalizeDiagnostics(ctx, &startingTimestamp, durationMetric, tags)

	if err := m._transport.SendMessage(ctx, request); err != nil {
		addExceptionTags(&tags, err)
		return nil, err
	}

	tasks.RegisterCancellation(ctx, func() {
		data, err := json.Marshal(messages.CancelledNotification{RequestId: *request.Id})
		if err != nil {
			return
		}
		_ = m.SendMessage(ctx, messages.NewJsonRpcNotification(
			messages.NotificationMethods_CancelledNotification,
			data,
		))
	})

	response, err := tcs.Result()
	if err != nil {
		addExceptionTags(&tags, err)
		return nil, err
	}

	switch resp := response.(type) {
	case *messages.JsonRpcError:
		return nil, fmt.Errorf("request failed (server side): %s", resp.Error.Message)
	case *messages.JsonRpcResponse:
		return resp, nil
	}

	return nil, fmt.Errorf("invalid response type")
}

func (m *McpSession) SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error {
	if m == nil || message == nil {
		return fmt.Errorf("mcp session or message is nil")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	durationMetric := s_clientOperationDuration
	if m._isServer {
		durationMetric = s_serverOperationDuration
	}

	method := getMethodName(message)

	var startingTimestamp int64 = time.Now().UnixNano()
	ctx, span := mcp.Tracer.Start(ctx, m.createActivityName(method))
	defer span.End()

	mcp.PropagatorInject(ctx, message)

	tags := []attribute.KeyValue{}
	m.addStandardTags(&tags, method)
	defer finalizeDiagnostics(ctx, &startingTimestamp, durationMetric, tags)

	if err := m._transport.SendMessage(ctx, message); err != nil {
		addExceptionTags(&tags, err)
		return err
	}

	if notification, ok := message.(*messages.JsonRpcNotification); ok {
		if params := getCancelledNotificationParams(notification.Params); params != nil {
			if c, ok := m._pendingRequests.Load(params.RequestId); ok {
				if source, ok := c.(*tasks.TaskCompletionSource[messages.IJsonRpcMessage]); ok {
					source.Cancel()
					m._pendingRequests.Delete(params.RequestId)
				}
			}
		}
	}

	return nil
}

func (m *McpSession) Dispose() {
	durationMetric := s_clientSessionDuration
	if m._isServer {
		durationMetric = s_serverSessionDuration
	}

	tags := []attribute.KeyValue{}
	tags = append(tags, attribute.String("session.id", m._id))
	tags = append(tags, attribute.String("network.transport", m._transportKind))

	incr := m._sessionStartingTimestamp - time.Now().UnixNano()
	durationMetric.Record(context.Background(), (float64)(incr), metric.WithAttributes(tags...))

	m._pendingRequests.Range(func(key, value interface{}) bool {
		if tcs, ok := value.(*tasks.TaskCompletionSource[messages.IJsonRpcMessage]); ok {
			tcs.Cancel()
		}
		return true
	})

	m._pendingRequests.Range(func(key, value any) bool {
		m._pendingRequests.Delete(key)
		return true
	})
}

func getCancelledNotificationParams(notificationParams interface{}) *messages.CancelledNotification {
	d, err := json.Marshal(notificationParams)
	if err != nil {
		return nil
	}
	var p messages.CancelledNotification
	err = json.Unmarshal(d, &p)
	if err != nil {
		return nil
	}
	return &p
}

func getMethodName(message messages.IJsonRpcMessage) string {
	switch request := message.(type) {
	case *messages.JsonRpcRequest:
		return request.Method
	case *messages.JsonRpcNotification:
		return request.Method
	default:
		return "unknownMethod"
	}
}

func (m *McpSession) addStandardTags(tags *[]attribute.KeyValue, method string) {
	*tags = append(*tags, attribute.String("session.id", m._id))
	*tags = append(*tags, attribute.String("rpc.system", "jsonrpc"))
	*tags = append(*tags, attribute.String("rpc.jsonrpc.version", "2.0"))
	*tags = append(*tags, attribute.String("rpc.method", method))
	*tags = append(*tags, attribute.String("etwork.transport", m._transportKind))
}

func addExceptionTags(tags *[]attribute.KeyValue, err error) {
	*tags = append(*tags, attribute.String("error", err.Error()))
	*tags = append(*tags, attribute.String("rpc.jsonrpc.error_code", "500")) //TODO: get error code from jsonrpc error
}

func finalizeDiagnostics(ctx context.Context, startingTimestamp *int64, durationMetric metric.Float64Histogram, tags []attribute.KeyValue) {
	if startingTimestamp != nil {
		incr := *startingTimestamp - time.Now().UnixNano()
		durationMetric.Record(ctx, (float64)(incr), metric.WithAttributes(tags...))
	}
}

func addRpcRequestTags(tags *[]attribute.KeyValue, request messages.JsonRpcRequest) {
	*tags = append(*tags, attribute.String("rpc.jsonrpc.request_id", request.Id.String()))
	if request.Params != nil {
		d, err := json.Marshal(request.Params)
		if err != nil {
			return
		}
		var p map[string]interface{}
		err = json.Unmarshal(d, &p)
		if err != nil {
			return
		}

		switch request.Method {
		case messages.RequestMethods_ToolsCall, messages.RequestMethods_PromptsGet:
			if prop, ok := p["name"].(string); ok {
				*tags = append(*tags, attribute.String("mcp.request.params.name", prop))
			}
		case messages.RequestMethods_ResourcesRead:
			if prop, ok := p["uri"].(string); ok {
				*tags = append(*tags, attribute.String("mcp.request.params.uri", prop))
			}

		}
	}
}
