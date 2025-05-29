package client

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type IClientTransport interface {
	GetName() string
	Connect(context.Context) (protocol.ITransport, error)
}
