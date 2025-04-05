package mcp

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

type IMcpEndpoint interface {
	SendRequest(ctx context.Context, request messages.JsonRpcRequest) (*messages.JsonRpcResponse, error)
	SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error
}
