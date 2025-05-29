package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type McpServerHandlers struct {
	ListToolsHandler   func(context.Context, RequestContext[*protocol.ListToolsRequestParams]) (*protocol.ListToolsResult, error)
	CallToolHandler    func(context.Context, RequestContext[*protocol.CallToolRequestParams]) (*protocol.CallToolResult, error)
	ListPromptsHandler func(context.Context, RequestContext[*protocol.ListPromptsRequestParams]) (*protocol.ListPromptsResult, error)
	GetPromptHandler   func(context.Context, RequestContext[*protocol.GetPromptRequestParams]) (*protocol.GetPromptResult, error)

	ListResourceTemplatesHandler    func(ctx context.Context, req RequestContext[*protocol.ListResourceTemplatesRequestParams]) (*protocol.ListResourceTemplatesResult, error) `json:"-"`
	ListResourcesHandler            func(ctx context.Context, req RequestContext[*protocol.ListResourcesRequestParams]) (*protocol.ListResourcesResult, error)                 `json:"-"`
	ReadResourceHandler             func(ctx context.Context, req RequestContext[*protocol.ReadResourceRequestParams]) (*protocol.ReadResourceResult, error)                   `json:"-"`
	SubscribeToResourcesHandler     func(ctx context.Context, req RequestContext[*protocol.SubscribeRequestParams]) (*protocol.EmptyResult, error)                             `json:"-"`
	UnsubscribeFromResourcesHandler func(ctx context.Context, req RequestContext[*protocol.UnsubscribeRequestParams]) (*protocol.EmptyResult, error)

	CompleteHandler        func(context.Context, RequestContext[*protocol.CompleteRequestParams]) (*protocol.CompleteResult, error)
	SetLoggingLevelHandler func(context.Context, RequestContext[*protocol.SetLevelRequestParams]) (*protocol.EmptyResult, error)
}

func (h *McpServerHandlers) OverwriteWithSetHandlers(option *McpServerOptions) {
	if option.Capabilities == nil {
		option.Capabilities = &ServerCapabilities{}
	}
	promptsCapability := option.Capabilities.Prompts

	if h.ListPromptsHandler != nil || h.GetPromptHandler != nil {
		if promptsCapability == nil {
			promptsCapability = &PromptsCapability{}
		}
		if h.ListPromptsHandler != nil {
			promptsCapability.ListPromptsHandler = h.ListPromptsHandler
		}
		if h.GetPromptHandler != nil {
			promptsCapability.GetPromptHandler = h.GetPromptHandler
		}
	}

	resourcesCapability := option.Capabilities.Resources
	if h.ListResourcesHandler != nil || h.ReadResourceHandler != nil {
		if resourcesCapability == nil {
			resourcesCapability = &ResourcesCapability{}
		}
		if h.ListResourceTemplatesHandler != nil {
			resourcesCapability.ListResourceTemplatesHandler = h.ListResourceTemplatesHandler
		}
		if h.ListResourcesHandler != nil {
			resourcesCapability.ListResourcesHandler = h.ListResourcesHandler
		}
		if h.ReadResourceHandler != nil {
			resourcesCapability.ReadResourceHandler = h.ReadResourceHandler
		}

		if h.SubscribeToResourcesHandler != nil || h.UnsubscribeFromResourcesHandler != nil {
			if h.SubscribeToResourcesHandler != nil {
				resourcesCapability.SubscribeToResourcesHandler = h.SubscribeToResourcesHandler
			}
			if h.UnsubscribeFromResourcesHandler != nil {
				resourcesCapability.UnsubscribeFromResourcesHandler = h.UnsubscribeFromResourcesHandler
			}
		}
	}

	toolsCapability := option.Capabilities.Tools
	if h.ListToolsHandler != nil || h.CallToolHandler != nil {
		if toolsCapability == nil {
			toolsCapability = &ToolsCapability{}
		}
		if h.ListToolsHandler != nil {
			toolsCapability.ListToolsHandler = h.ListToolsHandler
		}
		if h.CallToolHandler != nil {
			toolsCapability.CallToolHandler = h.CallToolHandler
		}
	}

	loggingCapability := option.Capabilities.Logging
	if h.SetLoggingLevelHandler != nil {
		if loggingCapability == nil {
			loggingCapability = &LoggingCapability{}
		}
		if h.SetLoggingLevelHandler != nil {
			loggingCapability.SetLoggingLevelHandler = h.SetLoggingLevelHandler
		}
	}

	completionsCapability := option.Capabilities.Completions
	if h.CompleteHandler != nil {
		if completionsCapability == nil {
			completionsCapability = &CompletionsCapability{}
		}
		if h.CompleteHandler != nil {
			completionsCapability.CompleteHandler = h.CompleteHandler
		}
	}

	option.Capabilities.Prompts = promptsCapability
	option.Capabilities.Resources = resourcesCapability
	option.Capabilities.Tools = toolsCapability
	option.Capabilities.Logging = loggingCapability
	option.Capabilities.Completions = completionsCapability
}
