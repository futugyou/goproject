package shared

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

type IMcpEndpoint interface {
	SendRequest(ctx context.Context, req *messages.JsonRpcRequest) (*messages.JsonRpcResponse, error)
	SendMessage(ctx context.Context, msg messages.IJsonRpcMessage) error
	RegisterNotificationHandler(method string, handler NotificationHandler) *RegistrationHandle
	GetEndpointName() string
	GetMessageProcessingTask() <-chan struct{}
	Dispose(ctx context.Context) error
	SendNotification(ctx context.Context, notification messages.JsonRpcNotification) error
	NotifyProgress(ctx context.Context, progressToken messages.ProgressToken, progress messages.ProgressNotificationValue) error
}

type Disposable interface {
	Dispose() error
}
