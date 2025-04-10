package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

// ResourcesCapability represents the resources capability configuration.
type ResourcesCapability struct {
	Subscribe                       *bool                                                                                                                                `json:"subscribe,omitempty"`
	ListChanged                     *bool                                                                                                                                `json:"listChanged,omitempty"`
	ListResourceTemplatesHandler    func(ctx context.Context, req RequestContext[*types.ListResourceTemplatesRequestParams]) (*types.ListResourceTemplatesResult, error) `json:"-"`
	ListResourcesHandler            func(ctx context.Context, req RequestContext[*types.ListResourcesRequestParams]) (*types.ListResourcesResult, error)                 `json:"-"`
	ReadResourceHandler             func(ctx context.Context, req RequestContext[*types.ReadResourceRequestParams]) (*types.ReadResourceResult, error)                   `json:"-"`
	SubscribeToResourcesHandler     func(ctx context.Context, req RequestContext[*types.SubscribeRequestParams]) (*types.EmptyResult, error)                             `json:"-"`
	UnsubscribeFromResourcesHandler func(ctx context.Context, req RequestContext[*types.UnsubscribeRequestParams]) (*types.EmptyResult, error)                           `json:"-"`
}

// LoggingCapability represents the logging capability configuration.
type LoggingCapability struct {
	SetLoggingLevelHandler func(ctx context.Context, req RequestContext[*types.SetLevelRequestParams]) (*types.EmptyResult, error) `json:"-"`
}

// ToolsCapability represents the tools capability configuration.
type ToolsCapability struct {
	ListChanged      *bool                                                                                                        `json:"listChanged,omitempty"`
	ListToolsHandler func(ctx context.Context, req RequestContext[*types.ListToolsRequestParams]) (*types.ListToolsResult, error) `json:"-"`
	CallToolHandler  func(ctx context.Context, req RequestContext[*types.CallToolRequestParams]) (*types.CallToolResult, error)   `json:"-"`
	ToolCollection   McpServerPrimitiveCollection[IMcpServerTool]                                                                 `json:"-"`
}
