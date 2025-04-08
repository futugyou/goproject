package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type McpServerPrompt struct {
	ProtocolPrompt types.Prompt
}

// GetName implements IMcpServerPrimitive.
func (m *McpServerPrompt) GetName() string {
	return m.ProtocolPrompt.Name
}

func (m *McpServerPrompt) GetProtocolPrompt() types.Prompt {
	return m.ProtocolPrompt
}

func (m *McpServerPrompt) Get(context.Context, RequestContext[*types.GetPromptRequestParams]) (*types.GetPromptResult, error) {
	panic("implement me")
}

var _ IMcpServerPrimitive = (*McpServerPrompt)(nil)
