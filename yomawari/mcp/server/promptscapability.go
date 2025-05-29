package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

// PromptsCapability represents the prompts capability configuration.
type PromptsCapability struct {
	ListChanged        *bool                                                                                                          `json:"listChanged,omitempty"`
	PromptCollection   *McpServerPrimitiveCollection[IMcpServerPrompt]                                                                `json:"-"`
	ListPromptsHandler func(context.Context, RequestContext[*protocol.ListPromptsRequestParams]) (*protocol.ListPromptsResult, error) `json:"-"`
	GetPromptHandler   func(context.Context, RequestContext[*protocol.GetPromptRequestParams]) (*protocol.GetPromptResult, error)     `json:"-"`
}
