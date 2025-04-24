package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/futugyou/yomawari/extensions-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/mcp/protocol/types"
	"github.com/futugyou/yomawari/mcp/shared"
)

var McpServerDefaultImplementation types.Implementation = types.Implementation{
	Name:    "McpServer",
	Version: "1.0.0",
}

var _ IMcpServer = (*McpServer)(nil)

type McpServer struct {
	*shared.BaseMcpEndpoint
	_sessionTransport        transport.ITransport
	_servicesScopePerRequest bool
	_toolsChangedDelegate    func()
	_promptsChangedDelegate  func()
	EndpointName             string
	_started                 int32
	ServerCapabilities       *ServerCapabilities
	ClientCapabilities       *shared.ClientCapabilities
	ClientInfo               *types.Implementation
	ServerOptions            McpServerOptions
}

func NewMcpServer(itransport transport.ITransport, options McpServerOptions) *McpServer {
	serverName := McpServerDefaultImplementation.Name
	version := McpServerDefaultImplementation.Version
	if len(options.ServerInfo.Name) > 0 {
		serverName = options.ServerInfo.Name
	}

	if len(options.ServerInfo.Version) > 0 {
		version = options.ServerInfo.Version
	}

	s := &McpServer{
		BaseMcpEndpoint:          shared.NewBaseMcpEndpoint(),
		_sessionTransport:        itransport,
		_servicesScopePerRequest: options.ScopeRequests,
		EndpointName:             fmt.Sprintf("Server (%s %s)", serverName, version),
		ServerOptions:            options,
	}

	s.setInitializeHandler(&options)
	s.setToolsHandler(&options)
	s.setPromptsHandler(&options)
	s.setResourcesHandler(&options)
	s.setCompletionHandler(&options)
	s.setPingHandler()

	if options.Capabilities != nil && len(options.Capabilities.NotificationHandlers) > 0 {
		s.GetNotificationHandlers().RegisterRange(options.Capabilities.NotificationHandlers)
	}

	if options.Capabilities != nil && options.Capabilities.Tools != nil && options.Capabilities.Tools.ToolCollection.Count() > 0 {
		s._toolsChangedDelegate = func() {
			s.SendMessage(context.Background(), transport.NewJsonRpcNotification(transport.NotificationMethods_ToolListChangedNotification, nil))
		}
		options.Capabilities.Tools.ToolCollection.OnChanged(s._toolsChangedDelegate)
	}

	if options.Capabilities != nil && options.Capabilities.Prompts != nil && options.Capabilities.Prompts.PromptCollection.Count() > 0 {
		s._promptsChangedDelegate = func() {
			s.SendMessage(context.Background(), transport.NewJsonRpcNotification(transport.NotificationMethods_PromptListChangedNotification, nil))
		}
		options.Capabilities.Prompts.PromptCollection.OnChanged(s._promptsChangedDelegate)
	}

	s.InitializeSession(itransport, true)
	return s
}

func (m *McpServer) setInitializeHandler(options *McpServerOptions) {
	shared.GenericRequestHandlerAdd(
		m.GetRequestHandlers(),
		transport.RequestMethods_Initialize,
		func(ctx context.Context, request *shared.InitializeRequestParams, tran transport.ITransport) (*InitializeResult, error) {
			m.ClientCapabilities = request.Capabilities
			m.ClientInfo = &request.ClientInfo
			_endpointName := fmt.Sprintf("%s, Client (%s %s)", m.EndpointName, m.ClientInfo.Name, m.ClientInfo.Version)
			m.GetMcpSession().EndpointName = _endpointName
			return &InitializeResult{
				ProtocolVersion: options.ProtocolVersion,
				Capabilities:    *m.ServerCapabilities,
				ServerInfo:      options.ServerInfo,
				Instructions:    options.ServerInstructions,
			}, nil
		},
		nil,
		nil,
	)
}

func (m *McpServer) setToolsHandler(options *McpServerOptions) {
	var toolsCapability *ToolsCapability
	var listToolsHandler func(ctx context.Context, req RequestContext[*types.ListToolsRequestParams]) (*types.ListToolsResult, error)
	var callToolHandler func(ctx context.Context, req RequestContext[*types.CallToolRequestParams]) (*types.CallToolResult, error)
	var tools *McpServerPrimitiveCollection[IMcpServerTool]
	if options.Capabilities != nil {
		toolsCapability = options.Capabilities.Tools
	}
	if toolsCapability != nil {
		listToolsHandler = toolsCapability.ListToolsHandler
		callToolHandler = toolsCapability.CallToolHandler
		tools = toolsCapability.ToolCollection
	}

	if tools != nil && !tools.IsEmpty() {
		originalListToolsHandler := listToolsHandler
		listToolsHandler = func(ctx context.Context, req RequestContext[*types.ListToolsRequestParams]) (*types.ListToolsResult, error) {
			result := &types.ListToolsResult{
				Tools: make([]types.Tool, 0, tools.Count()),
			}
			if originalListToolsHandler != nil {
				r, err := originalListToolsHandler(ctx, req)
				if err != nil {
					result = r
				}
			}

			if req.Params != nil && req.Params.Cursor != nil {
				for _, tool := range tools.ToSlice() {
					result.Tools = append(result.Tools, *tool.GetProtocolTool())
				}
			}

			return result, nil
		}

		originalCallToolHandler := callToolHandler
		callToolHandler = func(ctx context.Context, req RequestContext[*types.CallToolRequestParams]) (*types.CallToolResult, error) {
			var tool IMcpServerTool
			var ok bool
			if req.Params != nil {
				if tool, ok = tools.Get(req.Params.Name); !ok {
					if originalCallToolHandler != nil {
						return originalCallToolHandler(ctx, req)
					}
				}
			}

			return tool.Invoke(ctx, req)
		}
		listChanged := true
		m.ServerCapabilities = &ServerCapabilities{
			Experimental: options.Capabilities.Experimental,
			Logging:      options.Capabilities.Logging,
			Prompts:      options.Capabilities.Prompts,
			Resources:    options.Capabilities.Resources,
			Tools: &ToolsCapability{
				ListChanged:      &listChanged,
				ListToolsHandler: listToolsHandler,
				CallToolHandler:  callToolHandler,
				ToolCollection:   tools,
			},
			Completions:          &CompletionsCapability{},
			NotificationHandlers: map[string]shared.NotificationHandler{},
		}
	} else {
		m.ServerCapabilities = options.Capabilities
	}

	setHandler(m, transport.RequestMethods_ToolsList, listToolsHandler, nil, nil)
	setHandler(m, transport.RequestMethods_ToolsCall, callToolHandler, nil, nil)
}

func (m *McpServer) setPromptsHandler(options *McpServerOptions) {
	var promptsCapability *PromptsCapability
	var listPromptsHandler func(context.Context, RequestContext[*types.ListPromptsRequestParams]) (*types.ListPromptsResult, error)
	var getPromptHandler func(context.Context, RequestContext[*types.GetPromptRequestParams]) (*types.GetPromptResult, error)
	var prompts *McpServerPrimitiveCollection[IMcpServerPrompt]
	if options.Capabilities != nil {
		promptsCapability = options.Capabilities.Prompts
	}
	if promptsCapability != nil {
		listPromptsHandler = promptsCapability.ListPromptsHandler
		getPromptHandler = promptsCapability.GetPromptHandler
		prompts = promptsCapability.PromptCollection
	}
	if (listPromptsHandler == nil) != (getPromptHandler == nil) {
		return
	}
	if prompts != nil && !prompts.IsEmpty() {
		originalListPromptsHandler := listPromptsHandler
		listPromptsHandler = func(ctx context.Context, request RequestContext[*types.ListPromptsRequestParams]) (*types.ListPromptsResult, error) {
			result := &types.ListPromptsResult{
				Prompts: make([]types.Prompt, 0, prompts.Count()),
			}
			if originalListPromptsHandler != nil {
				r, err := originalListPromptsHandler(ctx, request)
				if err != nil {
					result = r
				}
			}

			if request.Params != nil && request.Params.Cursor != nil {
				for _, prompt := range prompts.ToSlice() {
					result.Prompts = append(result.Prompts, *prompt.GetProtocolPrompt())
				}
			}

			return result, nil
		}

		originalGetPromptHandler := getPromptHandler
		getPromptHandler = func(ctx context.Context, request RequestContext[*types.GetPromptRequestParams]) (*types.GetPromptResult, error) {
			var prompt IMcpServerPrompt
			var ok bool
			if request.Params != nil {
				if prompt, ok = prompts.Get(request.Params.Name); !ok {
					if originalGetPromptHandler != nil {
						return originalGetPromptHandler(ctx, request)
					}
				}
			}

			return prompt.Get(ctx, request)
		}
		listChanged := true
		m.ServerCapabilities = &ServerCapabilities{
			Experimental:         options.Capabilities.Experimental,
			Logging:              options.Capabilities.Logging,
			Resources:            options.Capabilities.Resources,
			Tools:                options.Capabilities.Tools,
			Completions:          &CompletionsCapability{},
			NotificationHandlers: map[string]shared.NotificationHandler{},
			Prompts: &PromptsCapability{
				ListChanged:        &listChanged,
				ListPromptsHandler: listPromptsHandler,
				GetPromptHandler:   getPromptHandler,
				PromptCollection:   prompts,
			},
		}
	} else {
		m.ServerCapabilities = options.Capabilities
	}

	setHandler(m, transport.RequestMethods_PromptsList, listPromptsHandler, nil, nil)
	setHandler(m, transport.RequestMethods_PromptsGet, getPromptHandler, nil, nil)
}

func (m *McpServer) setResourcesHandler(options *McpServerOptions) {
	if options.Capabilities == nil || options.Capabilities.Resources == nil {
		return
	}
	resourcesCapability := options.Capabilities.Resources
	listResourcesHandler := resourcesCapability.ListResourcesHandler
	listResourceTemplatesHandler := resourcesCapability.ListResourceTemplatesHandler
	readResourceHandler := resourcesCapability.ReadResourceHandler
	if listResourcesHandler == nil {
		listResourcesHandler = func(context.Context, RequestContext[*types.ListResourcesRequestParams]) (*types.ListResourcesResult, error) {
			return &types.ListResourcesResult{}, nil
		}
	}
	if listResourceTemplatesHandler == nil {
		listResourceTemplatesHandler = func(context.Context, RequestContext[*types.ListResourceTemplatesRequestParams]) (*types.ListResourceTemplatesResult, error) {
			return &types.ListResourceTemplatesResult{}, nil
		}
	}

	setHandler(m, transport.RequestMethods_ResourcesList, listResourcesHandler, nil, nil)
	setHandler(m, transport.RequestMethods_ResourcesRead, readResourceHandler, nil, nil)
	setHandler(m, transport.RequestMethods_ResourcesTemplatesList, listResourceTemplatesHandler, nil, nil)

	if resourcesCapability.Subscribe == nil && !*resourcesCapability.Subscribe {
		return
	}

	subscribeHandler := resourcesCapability.SubscribeToResourcesHandler
	unsubscribeHandler := resourcesCapability.UnsubscribeFromResourcesHandler
	if subscribeHandler == nil || unsubscribeHandler == nil {
		return
	}

	setHandler(m, transport.RequestMethods_ResourcesSubscribe, subscribeHandler, nil, nil)
	setHandler(m, transport.RequestMethods_ResourcesUnsubscribe, unsubscribeHandler, nil, nil)
}

func (m *McpServer) setCompletionHandler(options *McpServerOptions) {
	if options.Capabilities == nil || options.Capabilities.Completions == nil {
		return
	}
	completionsCapability := options.Capabilities.Completions
	completeHandler := completionsCapability.CompleteHandler
	if completeHandler == nil {
		return
	}

	setHandler(m, transport.RequestMethods_CompletionComplete, completeHandler, nil, nil)
}

func (m *McpServer) setPingHandler() {
	setHandler(
		m,
		transport.RequestMethods_Ping,
		func(context.Context, RequestContext[*types.PingRequest]) (*types.PingResult, error) {
			return &types.PingResult{}, nil
		}, nil, nil)
}

// GetClientCapabilities implements IMcpServer.
func (m *McpServer) GetClientCapabilities() *shared.ClientCapabilities {
	return m.ClientCapabilities
}

// GetClientInfo implements IMcpServer.
func (m *McpServer) GetClientInfo() *types.Implementation {
	return m.ClientInfo
}

// GetMcpServerOptions implements IMcpServer.
func (m *McpServer) GetMcpServerOptions() *McpServerOptions {
	return &m.ServerOptions
}

// Run implements IMcpServer.
func (m *McpServer) Run(ctx context.Context) error {
	if atomic.SwapInt32(&m._started, 1) != 0 {
		return fmt.Errorf("server already started")
	}

	m.StartSession(ctx, m._sessionTransport)
	<-m.GetMessageProcessingTask()

	return m.Dispose(ctx)
}

func (e *McpServer) AsSamplingChatClient() (chatcompletion.IChatClient, error) {
	if e.GetClientCapabilities() == nil || e.GetClientCapabilities().Sampling == nil {
		return nil, fmt.Errorf("client capabilities sampling not set")
	}
	return NewSamplingChatClient(e), nil
}

func (e *McpServer) RequestRoots(ctx context.Context, request types.ListRootsRequestParams) (*types.ListRootsResult, error) {
	if e.GetClientCapabilities() == nil || e.GetClientCapabilities().Roots == nil {
		return nil, fmt.Errorf("client capabilities roots not set")
	}
	req := transport.NewJsonRpcRequest(transport.RequestMethods_RootsList, request, nil)
	resp, err := e.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	var result types.ListRootsResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func setHandler[TRequest any, TResponse any](
	m *McpServer,
	method string,
	handler func(context.Context, RequestContext[*TRequest]) (*TResponse, error),
	unmarshaler shared.RequestUnmarshaler[TRequest],
	marshaler shared.RepsonseMarshaler[TResponse],
) {
	shared.GenericRequestHandlerAdd(
		m.GetRequestHandlers(),
		method,
		func(ctx context.Context, request *TRequest, destinationTransport transport.ITransport) (*TResponse, error) {
			return invokeHandler(m, ctx, handler, request, destinationTransport)
		},
		unmarshaler,
		marshaler,
	)
}

func invokeHandler[TParams any, TResult any](
	m *McpServer,
	ctx context.Context,
	handler func(context.Context, RequestContext[TParams]) (*TResult, error),
	args TParams,
	destinationTransport transport.ITransport,
) (*TResult, error) {
	// TODO: handle _servicesScopePerRequest
	return handler(ctx, RequestContext[TParams]{Params: args})
}
