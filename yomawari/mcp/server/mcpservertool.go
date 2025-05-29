package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type IMcpServerTool interface {
	IMcpServerPrimitive
	GetProtocolTool() *protocol.Tool
	Invoke(ctx context.Context, request RequestContext[*protocol.CallToolRequestParams]) (*protocol.CallToolResult, error)
}
