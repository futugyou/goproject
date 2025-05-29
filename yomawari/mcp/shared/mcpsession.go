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
	"github.com/futugyou/yomawari/mcp/protocol"
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
	_transport                protocol.ITransport
	_requestHandlers          *RequestHandlers
	_notificationHandlers     *NotificationHandlers
	_sessionStartingTimestamp int64
	_pendingRequests          sync.Map // map[RequestId]*tasks.TaskCompletionSource[protocol.IJsonRpcMessage]
	_handlingRequests         sync.Map // map[RequestId]context.CancelFunc
	_id                       string
	_nextRequestId            int64
	EndpointName              string
}

func NewMcpSession(isServer bool, transp protocol.ITransport, endpointName string, requestHandlers *RequestHandlers, notificationHandlers *NotificationHandlers) *McpSession {
	transportKind := transp.GetTransportKind()

	return &McpSession{
		_isServer:                 isServer,
		_transportKind:            string(transportKind),
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
	var processMessage = func(ctx context.Context, message protocol.IJsonRpcMessage) error {
		var messageWithId protocol.IJsonRpcMessageWithId
		if msg, ok := message.(protocol.IJsonRpcMessageWithId); ok {
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
			if request, ok := message.(*protocol.JsonRpcRequest); !isUserCancellation && ok {
				m._transport.SendMessage(ctx, protocol.NewJsonRpcErrorWithTransport(request.Id, 500, err.Error(), nil, request.RelatedTransport))
			}
		}
		return nil
	}

	resultChan := m._transport.MessageReader()
	processor := tasks.TaskProcessor[protocol.IJsonRpcMessage]{
		ResultChan:        resultChan,
		Handler:           processMessage, // func(ctx, msg) error
		MaxConcurrency:    20,
		PerMessageTimeout: 10 * time.Second,
	}

	processor.Run(ctx)

	m._pendingRequests.Range(func(key, value interface{}) bool {
		if tcs, ok := value.(*tasks.TaskCompletionSource[protocol.IJsonRpcMessage]); ok {
			tcs.TrySetError(fmt.Errorf("the server shut down unexpectedly"))
		}
		return true
	})

	return nil
}

func (m *McpSession) handleMessage(ctx context.Context, message protocol.IJsonRpcMessage) error {
	durationMetric := s_clientOperationDuration
	if m._isServer {
		durationMetric = s_serverOperationDuration
	}
	method := getMethodName(message)

	var startingTimestamp int64 = time.Now().UnixNano()
	ctx, span := StartSpanWithJsonRpcData(ctx, m.createActivityName(method), message)
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
	case *protocol.JsonRpcRequest:
		addRpcRequestTags(&tags, *request)
		err = m.handleRequest(ctx, *request)
	case *protocol.JsonRpcNotification:
		err = m.handleNotification(ctx, request)
	case protocol.IJsonRpcMessageWithId:
		err = m.handleMessageWithId(message, request)
	default:
	}

	if err != nil {
		addExceptionTags(&tags, err)
		return err
	}

	return nil
}

func (m *McpSession) handleNotification(ctx context.Context, notification *protocol.JsonRpcNotification) error {
	if notification.Method == protocol.NotificationMethods_CancelledNotification {
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

func (m *McpSession) handleMessageWithId(message protocol.IJsonRpcMessage, messageWithId protocol.IJsonRpcMessageWithId) error {
	if messageWithId.GetId() == nil || len(messageWithId.GetId().String()) == 0 {
		return fmt.Errorf("message with id has no id")
	}
	requestid := messageWithId.GetId()
	value, ok := m._pendingRequests.Load(*requestid)
	if !ok {
		return fmt.Errorf("no pending request found for id %s", requestid.String())
	}
	if source, ok := value.(*tasks.TaskCompletionSource[protocol.IJsonRpcMessage]); ok {
		source.SetResult(message)
	}
	return nil
}

func (m *McpSession) handleRequest(ctx context.Context, request protocol.JsonRpcRequest) error {
	handler, ok := m._requestHandlers.Get(request.Method)
	if !ok {
		return fmt.Errorf("no handler found for method %s", request.Method)
	}

	result, err := handler(ctx, &request)
	if err != nil {
		return err
	}

	msg := protocol.NewJsonRpcResponseWithTransport(request.Id, result, request.RelatedTransport)
	return m._transport.SendMessage(ctx, msg)
}

func (m *McpSession) RegisterNotificationHandler(method string, handler protocol.NotificationHandler) *RegistrationHandle {
	return m._notificationHandlers.Register(method, handler, true)
}

func (m *McpSession) createActivityName(method string) string {
	s := "client"
	if m._isServer {
		s = "server"
	}

	return fmt.Sprintf("mcp.%s.%s/%s", s, m._transportKind, method)
}

func (m *McpSession) SendRequest(ctx context.Context, request *protocol.JsonRpcRequest) (*protocol.JsonRpcResponse, error) {
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
		request.Id = protocol.NewRequestIdFromString(fmt.Sprintf("%s-%d", m._id, newId))
	}

	PropagatorInject(ctx, request)

	ctx, cancelfunc := context.WithCancel(ctx)
	tcs := tasks.NewTaskCompletionSource[protocol.IJsonRpcMessage](ctx, cancelfunc)
	m._pendingRequests.Store(request.Id, tcs)

	m.addStandardTags(&tags, method)
	addRpcRequestTags(&tags, *request)

	defer finalizeDiagnostics(ctx, &startingTimestamp, durationMetric, tags)

	m.sendToRelatedTransport(ctx, request)

	tasks.RegisterCancellation(ctx, func() {
		data, err := json.Marshal(protocol.CancelledNotification{RequestId: *request.Id})
		if err != nil {
			return
		}
		_ = m.SendMessage(ctx, protocol.NewJsonRpcNotificationWithTransport(
			protocol.NotificationMethods_CancelledNotification,
			data,
			request.RelatedTransport,
		))
	})

	response, err := tcs.Result()
	if err != nil {
		addExceptionTags(&tags, err)
		return nil, err
	}

	switch resp := response.(type) {
	case *protocol.JsonRpcError:
		return nil, fmt.Errorf("request failed (server side): %s", resp.Error.Message)
	case *protocol.JsonRpcResponse:
		return resp, nil
	}

	return nil, fmt.Errorf("invalid response type")
}

func (m *McpSession) SendMessage(ctx context.Context, message protocol.IJsonRpcMessage) error {
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

	PropagatorInject(ctx, message)

	tags := []attribute.KeyValue{}
	m.addStandardTags(&tags, method)
	defer finalizeDiagnostics(ctx, &startingTimestamp, durationMetric, tags)

	m.sendToRelatedTransport(ctx, message)

	if notification, ok := message.(*protocol.JsonRpcNotification); ok {
		if params := getCancelledNotificationParams(notification.Params); params != nil {
			if c, ok := m._pendingRequests.Load(params.RequestId); ok {
				if source, ok := c.(*tasks.TaskCompletionSource[protocol.IJsonRpcMessage]); ok {
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
		if tcs, ok := value.(*tasks.TaskCompletionSource[protocol.IJsonRpcMessage]); ok {
			tcs.Cancel()
		}
		return true
	})

	m._pendingRequests.Range(func(key, value any) bool {
		m._pendingRequests.Delete(key)
		return true
	})
}

func getCancelledNotificationParams(notificationParams interface{}) *protocol.CancelledNotification {
	d, err := json.Marshal(notificationParams)
	if err != nil {
		return nil
	}
	var p protocol.CancelledNotification
	err = json.Unmarshal(d, &p)
	if err != nil {
		return nil
	}
	return &p
}

func getMethodName(message protocol.IJsonRpcMessage) string {
	switch request := message.(type) {
	case *protocol.JsonRpcRequest:
		return request.Method
	case *protocol.JsonRpcNotification:
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

func (m *McpSession) sendToRelatedTransport(ctx context.Context, message protocol.IJsonRpcMessage) {
	transport := message.GetRelatedTransport()
	if transport == nil {
		transport = m._transport
	}
	transport.SendMessage(ctx, message)
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

func addRpcRequestTags(tags *[]attribute.KeyValue, request protocol.JsonRpcRequest) {
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
		case protocol.RequestMethods_ToolsCall, protocol.RequestMethods_PromptsGet:
			if prop, ok := p["name"].(string); ok {
				*tags = append(*tags, attribute.String("mcp.request.params.name", prop))
			}
		case protocol.RequestMethods_ResourcesRead:
			if prop, ok := p["uri"].(string); ok {
				*tags = append(*tags, attribute.String("mcp.request.params.uri", prop))
			}

		}
	}
}
