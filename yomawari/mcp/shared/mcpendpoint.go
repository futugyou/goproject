package shared

import (
	"context"
	"errors"
	"sync"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
	"github.com/futugyou/yomawari/mcp/protocol/transport"
)

var _ IMcpEndpoint = (*BaseMcpEndpoint)(nil)

type IMcpEndpoint interface {
	SendRequest(ctx context.Context, req *messages.JsonRpcRequest) (*messages.JsonRpcResponse, error)
	SendMessage(ctx context.Context, msg messages.IJsonRpcMessage) error
	RegisterNotificationHandler(method string, handler NotificationHandler) *RegistrationHandle
	GetEndpointName() string
	GetMessageProcessingTask() <-chan struct{}
	Dispose(ctx context.Context) error
}

type Disposable interface {
	Dispose() error
}

type BaseMcpEndpoint struct {
	mu            sync.Mutex
	disposed      bool
	session       *McpSession
	sessionCts    context.CancelFunc
	messageTask   <-chan struct{}
	reqHandlers   *RequestHandlers
	notifHandlers *NotificationHandlers
	endpointName  string
}

func NewBaseMcpEndpoint() *BaseMcpEndpoint {
	return &BaseMcpEndpoint{
		reqHandlers:   NewRequestHandlers(),
		notifHandlers: NewNotificationHandlers(),
		endpointName:  "",
	}
}

func (e *BaseMcpEndpoint) GetEndpointName() string {
	return e.endpointName
}

func (e *BaseMcpEndpoint) GetMessageProcessingTask() <-chan struct{} {
	return e.messageTask
}

func (e *BaseMcpEndpoint) GetRequestHandlers() *RequestHandlers {
	return e.reqHandlers
}

func (e *BaseMcpEndpoint) GetNotificationHandlers() *NotificationHandlers {
	return e.notifHandlers
}

func (e *BaseMcpEndpoint) InitializeSession(transport transport.ITransport, isServer bool) {
	e.session = NewMcpSession(isServer, transport, e.endpointName, e.reqHandlers, e.notifHandlers)
}

func (e *BaseMcpEndpoint) StartSession(ctx context.Context, transport transport.ITransport) {
	childCtx, cancel := context.WithCancel(ctx)
	e.sessionCts = cancel

	done := make(chan struct{})
	e.messageTask = done

	go func() {
		defer close(done)
		e.session.ProcessMessages(childCtx)
	}()
}

func (e *BaseMcpEndpoint) CancelSession() {
	if e != nil && e.sessionCts != nil {
		e.sessionCts()
	}
}

func (e *BaseMcpEndpoint) SendRequest(ctx context.Context, req *messages.JsonRpcRequest) (*messages.JsonRpcResponse, error) {
	if e == nil || e.session == nil {
		return nil, errors.New("session not initialized")
	}
	return e.session.SendRequest(ctx, req)
}

func (e *BaseMcpEndpoint) SendMessage(ctx context.Context, msg messages.IJsonRpcMessage) error {
	if e == nil || e.session == nil {
		return errors.New("session not initialized")
	}
	return e.session.SendMessage(ctx, msg)
}

func (e *BaseMcpEndpoint) RegisterNotificationHandler(method string, handler NotificationHandler) *RegistrationHandle {
	if e.session == nil {
		return nil
	}
	return e.session.RegisterNotificationHandler(method, handler)
}

func (e *BaseMcpEndpoint) Dispose(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.disposed {
		return nil
	}
	e.disposed = true
	return e.disposeUnsynchronized(ctx)
}

func (e *BaseMcpEndpoint) disposeUnsynchronized(ctx context.Context) error {
	if e.sessionCts != nil {
		e.sessionCts()
	}

	if e.messageTask != nil {
		select {
		case <-e.messageTask:
		case <-ctx.Done():
		}
	}

	e.session.Dispose()
	return nil
}
