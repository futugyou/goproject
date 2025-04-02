package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

type IMcpEndpoint interface {
	SendRequest(ctx context.Context, request messages.JsonRpcRequest) (interface{}, error)
	SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error
	AddNotificationHandler(method string, hander func(messages.JsonRpcNotification) error)
}
