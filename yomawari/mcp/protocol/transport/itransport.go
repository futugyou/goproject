package transport

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

type ITransport interface {
	IsConnected() bool
	MessageReader() <-chan messages.IJsonRpcMessage
	SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error
	Close() error
}
