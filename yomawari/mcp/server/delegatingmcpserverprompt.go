package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

var _ IMcpServerPrompt = (*DelegatingMcpServerPrompt)(nil)

type DelegatingMcpServerPrompt struct {
	innerPrompt IMcpServerPrompt
}

func NewDelegatingMcpServerPrompt(innerPrompt IMcpServerPrompt) *DelegatingMcpServerPrompt {
	return &DelegatingMcpServerPrompt{
		innerPrompt: innerPrompt,
	}
}

// Get implements IMcpServerPrompt.
func (d *DelegatingMcpServerPrompt) Get(ctx context.Context, request RequestContext[*types.GetPromptRequestParams]) (*types.GetPromptResult, error) {
	return d.innerPrompt.Get(ctx, request)
}

// GetId implements IMcpServerPrompt.
func (d *DelegatingMcpServerPrompt) GetId() string {
	return d.innerPrompt.GetId()
}

// GetProtocolPrompt implements IMcpServerPrompt.
func (d *DelegatingMcpServerPrompt) GetProtocolPrompt() *types.Prompt {
	return d.innerPrompt.GetProtocolPrompt()
}
