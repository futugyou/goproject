package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
	"github.com/futugyou/yomawari/mcp/shared"
)

type IMcpServer interface {
	shared.IMcpEndpoint
	GetClientCapabilities() *types.ClientCapabilities
	GetClientInfo() *types.Implementation
	GetMcpServerOptions() *McpServerOptions
	Run(ctx context.Context) error
}
