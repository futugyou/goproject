package server

import "github.com/futugyou/yomawari/mcp/protocol/types"

type McpServerPrompt struct {
	ProtocolPrompt types.Prompt
}

// GetName implements IMcpServerPrimitive.
func (m *McpServerPrompt) GetName() string {
	return m.ProtocolPrompt.Name
}

var _ IMcpServerPrimitive = (*McpServerPrompt)(nil)
