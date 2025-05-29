package shared

import (
	"context"
	
	"github.com/futugyou/yomawari/mcp/protocol"
)

type IMcpEndpoint interface {
	SendRequest(ctx context.Context, req *protocol.JsonRpcRequest) (*protocol.JsonRpcResponse, error)
	SendMessage(ctx context.Context, msg protocol.IJsonRpcMessage) error
	RegisterNotificationHandler(method string, handler protocol.NotificationHandler) *RegistrationHandle
	GetEndpointName() string
	GetMessageProcessingTask() <-chan struct{}
	Dispose(ctx context.Context) error
	SendNotification(ctx context.Context, notification protocol.JsonRpcNotification) error
	NotifyProgress(ctx context.Context, progressToken protocol.ProgressToken, progress protocol.ProgressNotificationValue) error
}

type Disposable interface {
	Dispose() error
}
