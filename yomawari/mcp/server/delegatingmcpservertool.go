package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

var _ IMcpServerTool = (*DelegatingMcpServerTool)(nil)

type DelegatingMcpServerTool struct {
	innerTool IMcpServerTool
}

func NewDelegatingMcpServerTool(innerTool IMcpServerTool) *DelegatingMcpServerTool {
	return &DelegatingMcpServerTool{
		innerTool: innerTool,
	}
}

// GetId implements IMcpServerTool.
func (d *DelegatingMcpServerTool) GetId() string {
	return d.innerTool.GetId()
}

// GetProtocolTool implements IMcpServerTool.
func (d *DelegatingMcpServerTool) GetProtocolTool() *protocol.Tool {
	return d.innerTool.GetProtocolTool()
}

// Invoke implements IMcpServerTool.
func (d *DelegatingMcpServerTool) Invoke(ctx context.Context, request RequestContext[*protocol.CallToolRequestParams]) (*protocol.CallToolResult, error) {
	return d.innerTool.Invoke(ctx, request)
}
