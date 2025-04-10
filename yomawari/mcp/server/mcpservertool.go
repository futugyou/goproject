package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type IMcpServerTool interface {
	IMcpServerPrimitive
	GetProtocolTool() *types.Tool
	Invoke(ctx context.Context, request RequestContext[*types.CallToolRequestParams]) (*types.CallToolResult, error)
}
