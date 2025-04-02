package server

type McpServerTool struct {
}

// GetName implements IMcpServerPrimitive.
func (m *McpServerTool) GetName() string {
	panic("unimplemented")
}

var _ IMcpServerPrimitive = (*McpServerTool)(nil)
