package mcp

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

type IMcpEndpoint interface {
	SendRequest(ctx context.Context, request messages.JsonRpcRequest) (interface{}, error)
	SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error
	RegisterNotificationHandler(method string, hander func(context.Context, messages.JsonRpcNotification) error)
	SendNotification(ctx context.Context, notification messages.JsonRpcNotification) error
	NotifyProgress(ctx context.Context, progressToken messages.ProgressToken, progress messages.ProgressNotificationValue) error
}
