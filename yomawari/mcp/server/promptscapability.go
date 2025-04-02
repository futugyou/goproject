package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

// PromptsCapability represents the prompts capability configuration.
type PromptsCapability struct {
	ListChanged        *bool                                                                                            `json:"listChanged,omitempty"`
	ListPromptsHandler func(ctx context.Context, req *types.ListPromptsRequestParams) (*types.ListPromptsResult, error) `json:"-"`
	GetPromptHandler   func(ctx context.Context, req *types.GetPromptRequestParams) (*types.GetPromptResult, error)     `json:"-"`
	PromptCollection   *McpServerPrimitiveCollection[*McpServerPrompt]                                                  `json:"-"`
}
