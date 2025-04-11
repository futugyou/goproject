package client

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type McpClientPrompt struct {
	prompt types.Prompt
	client IMcpClient
}

func NewMcpClientPrompt(prompt types.Prompt, client IMcpClient) *McpClientPrompt {
	return &McpClientPrompt{
		prompt: prompt,
		client: client,
	}
}

func (p *McpClientPrompt) GetPrompt() *types.Prompt {
	return &p.prompt
}

func (p *McpClientPrompt) GetName() string {
	return p.prompt.Name
}

func (p *McpClientPrompt) GetDescription() *string {
	return p.prompt.Description
}

func (p *McpClientPrompt) Get(ctx context.Context, arguments map[string]interface{}) (*types.GetPromptResult, error) {
	return p.client.GetPrompt(ctx, p.prompt.Name, arguments)
}
