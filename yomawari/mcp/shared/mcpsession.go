package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/futugyou/yomawari/mcp"
	"github.com/futugyou/yomawari/mcp/protocol/messages"
	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var s_clientSessionDuration = mcp.CreateDurationHistogram("mcp.client.session.duration", "Measures the duration of a client session.", true)
var s_serverSessionDuration = mcp.CreateDurationHistogram("mcp.server.session.duration", "Measures the duration of a server session.", true)
var s_clientRequestDuration = mcp.CreateDurationHistogram("rpc.client.duration", "Measures the duration of outbound RPC.", true)
var s_serverRequestDuration = mcp.CreateDurationHistogram("rpc.server.duration", "Measures the duration of inbound RPC.", true)

type McpSession struct {
	_isServer                 bool
	_transportKind            string
	_transport                transport.ITransport
	_requestHandlers          *RequestHandlers
	_notificationHandlers     *NotificationHandlers
	_sessionStartingTimestamp int64
	_pendingRequests          sync.Map // map[RequestId]*responseWrapper
	_handlingRequests         sync.Map
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

	if !m._transport.IsConnected() {
		return nil, fmt.Errorf("transport is not connected")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	durationMetric := s_clientRequestDuration
	if m._isServer {
		durationMetric = s_serverRequestDuration
	}

	method := request.Method

	var startingTimestamp int64 = time.Now().UnixNano()
	ctx, span := mcp.Tracer.Start(ctx, m.createActivityName(method))

	tags := []attribute.KeyValue{}
	m.addStandardTags(&tags, method)
	defer finalizeDiagnostics(ctx, &startingTimestamp, durationMetric, span, tags)

	// Set request ID
	if request.Id == nil {
		newId := atomic.AddInt64(&m._nextRequestId, 1)
		request.Id = messages.NewRequestIdFromString(fmt.Sprintf("%s-%d", m._id, newId))
	}

	ctx, cancelfunc := context.WithCancel(ctx)
	tcs := NewTaskCompletionSource[messages.IJsonRpcMessage](ctx, cancelfunc)
	m._pendingRequests.Store(request.Id, tcs)

	m.addStandardTags(&tags, method)
	addRpcRequestTags(&tags, *request)

	defer finalizeDiagnostics(ctx, &startingTimestamp, durationMetric, span, tags)

	if err := m._transport.SendMessage(ctx, request); err != nil {
		addExceptionTags(&tags, err)
		return nil, err
	}

	RegisterCancellation(ctx, func() {
		_ = m.SendMessage(ctx, messages.NewJsonRpcNotification(
			messages.NotificationMethods_CancelledNotification,
			messages.CancelledNotification{RequestId: *request.Id},
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

	if !m._transport.IsConnected() {
		return fmt.Errorf("transport is not connected")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	durationMetric := s_clientRequestDuration
	if m._isServer {
		durationMetric = s_serverRequestDuration
	}

	method := getMethodName(message)

	var startingTimestamp int64 = time.Now().UnixNano()
	ctx, span := mcp.Tracer.Start(ctx, m.createActivityName(method))

	tags := []attribute.KeyValue{}
	m.addStandardTags(&tags, method)
	defer finalizeDiagnostics(ctx, &startingTimestamp, durationMetric, span, tags)

	if err := m._transport.SendMessage(ctx, message); err != nil {
		addExceptionTags(&tags, err)
		return err
	}

	if notification, ok := message.(*messages.JsonRpcNotification); ok {
		if params := getCancelledNotificationParams(notification.Params); params != nil {
			if c, ok := m._pendingRequests.Load(params.RequestId); ok {
				if cancel, ok := c.(context.CancelFunc); ok {
					cancel()
					m._pendingRequests.Delete(params.RequestId)
				}
			}
		}
	}

	return nil
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

func finalizeDiagnostics(ctx context.Context, startingTimestamp *int64, durationMetric metric.Float64Histogram, traceSpan trace.Span, tags []attribute.KeyValue) {
	if startingTimestamp != nil {
		incr := *startingTimestamp - time.Now().UnixNano()
		durationMetric.Record(ctx, (float64)(incr), metric.WithAttributes(tags...))
	}
	traceSpan.End()
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
