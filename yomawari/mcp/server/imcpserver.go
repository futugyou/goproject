package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp"
	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type IMcpServer interface {
	mcp.IMcpEndpoint
	GetClientCapabilities() types.ClientCapabilities
	GetClientInfo() types.Implementation
	Run(ctx context.Context) error
}
