package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type McpServerHandlers struct {
	ListToolsHandler   func(context.Context, RequestContext[*types.ListToolsRequestParams]) (*types.ListToolsResult, error)
	CallToolHandler    func(context.Context, RequestContext[*types.CallToolRequestParams]) (*types.CallToolResponse, error)
	ListPromptsHandler func(context.Context, RequestContext[*types.ListPromptsRequestParams]) (*types.ListPromptsResult, error)
	GetPromptHandler   func(context.Context, RequestContext[*types.GetPromptRequestParams]) (*types.GetPromptResult, error)

	ListResourceTemplatesHandler    func(ctx context.Context, req RequestContext[*types.ListResourceTemplatesRequestParams]) (*types.ListResourceTemplatesResult, error) `json:"-"`
	ListResourcesHandler            func(ctx context.Context, req RequestContext[*types.ListResourcesRequestParams]) (*types.ListResourcesResult, error)                 `json:"-"`
	ReadResourceHandler             func(ctx context.Context, req RequestContext[*types.ReadResourceRequestParams]) (*types.ReadResourceResult, error)                   `json:"-"`
	SubscribeToResourcesHandler     func(ctx context.Context, req RequestContext[*types.SubscribeRequestParams]) (*types.EmptyResult, error)                             `json:"-"`
	UnsubscribeFromResourcesHandler func(ctx context.Context, req RequestContext[*types.UnsubscribeRequestParams]) (*types.EmptyResult, error)

	CompleteHandler        func(context.Context, RequestContext[*types.CompleteRequestParams]) (*types.CompleteResult, error)
	SetLoggingLevelHandler func(context.Context, RequestContext[*types.SetLevelRequestParams]) (*types.EmptyResult, error)
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
