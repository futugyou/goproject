package shared

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type IMcpEndpoint interface {
	SendRequest(ctx context.Context, req *transport.JsonRpcRequest) (*transport.JsonRpcResponse, error)
	SendMessage(ctx context.Context, msg transport.IJsonRpcMessage) error
	RegisterNotificationHandler(method string, handler transport.NotificationHandler) *RegistrationHandle
	GetEndpointName() string
	GetMessageProcessingTask() <-chan struct{}
	Dispose(ctx context.Context) error
	SendNotification(ctx context.Context, notification transport.JsonRpcNotification) error
	NotifyProgress(ctx context.Context, progressToken types.ProgressToken, progress transport.ProgressNotificationValue) error
}

type Disposable interface {
	Dispose() error
}
