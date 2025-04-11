package client

import (
	"context"

	"github.com/futugyou/yomawari/mcp"
	"github.com/futugyou/yomawari/mcp/protocol/types"
	"github.com/futugyou/yomawari/mcp/server"
)

type IMcpClient interface {
	mcp.IMcpEndpoint
	GetServerCapabilities() *server.ServerCapabilities
	GetServerInfo() *types.Implementation
	GetServerInstructions() *string
	Ping(ctx context.Context) error
	ListTools(ctx context.Context) ([]McpClientTool, error)
	CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (*types.CallToolResult, error)
}
