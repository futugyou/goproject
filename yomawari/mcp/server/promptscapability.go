package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

// PromptsCapability represents the prompts capability configuration.
type PromptsCapability struct {
	ListChanged        *bool                                                                                                    `json:"listChanged,omitempty"`
	PromptCollection   *McpServerPrimitiveCollection[IMcpServerPrompt]                                                          `json:"-"`
	ListPromptsHandler func(context.Context, RequestContext[*types.ListPromptsRequestParams]) (*types.ListPromptsResult, error) `json:"-"`
	GetPromptHandler   func(context.Context, RequestContext[*types.GetPromptRequestParams]) (*types.GetPromptResult, error)     `json:"-"`
}
