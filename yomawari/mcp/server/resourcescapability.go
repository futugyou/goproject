package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

// ResourcesCapability represents the resources capability configuration.
type ResourcesCapability struct {
	Subscribe                       *bool                                                                                                                                      `json:"subscribe,omitempty"`
	ListChanged                     *bool                                                                                                                                      `json:"listChanged,omitempty"`
	ListResourceTemplatesHandler    func(ctx context.Context, req RequestContext[*protocol.ListResourceTemplatesRequestParams]) (*protocol.ListResourceTemplatesResult, error) `json:"-"`
	ListResourcesHandler            func(ctx context.Context, req RequestContext[*protocol.ListResourcesRequestParams]) (*protocol.ListResourcesResult, error)                 `json:"-"`
	ReadResourceHandler             func(ctx context.Context, req RequestContext[*protocol.ReadResourceRequestParams]) (*protocol.ReadResourceResult, error)                   `json:"-"`
	SubscribeToResourcesHandler     func(ctx context.Context, req RequestContext[*protocol.SubscribeRequestParams]) (*protocol.EmptyResult, error)                             `json:"-"`
	UnsubscribeFromResourcesHandler func(ctx context.Context, req RequestContext[*protocol.UnsubscribeRequestParams]) (*protocol.EmptyResult, error)                           `json:"-"`
	ResourceCollection              *McpServerPrimitiveCollection[IMcpServerResource]                                                                                          `json:"-"`
}

// LoggingCapability represents the logging capability configuration.
type LoggingCapability struct {
	SetLoggingLevelHandler func(ctx context.Context, req RequestContext[*protocol.SetLevelRequestParams]) (*protocol.EmptyResult, error) `json:"-"`
}

// ToolsCapability represents the tools capability configuration.
type ToolsCapability struct {
	ListChanged      *bool                                                                                                              `json:"listChanged,omitempty"`
	ListToolsHandler func(ctx context.Context, req RequestContext[*protocol.ListToolsRequestParams]) (*protocol.ListToolsResult, error) `json:"-"`
	CallToolHandler  func(ctx context.Context, req RequestContext[*protocol.CallToolRequestParams]) (*protocol.CallToolResult, error)   `json:"-"`
	ToolCollection   *McpServerPrimitiveCollection[IMcpServerTool]                                                                      `json:"-"`
}
