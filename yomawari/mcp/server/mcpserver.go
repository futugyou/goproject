package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/futugyou/yomawari/mcp"
	"github.com/futugyou/yomawari/mcp/protocol"
	"github.com/futugyou/yomawari/mcp/shared"
)

var McpServerDefaultImplementation protocol.Implementation = protocol.Implementation{
	Name:    "McpServer",
	Version: "1.0.0",
}

var _ IMcpServer = (*McpServer)(nil)

type McpServer struct {
	*shared.BaseMcpEndpoint
	_sessionTransport        protocol.ITransport
	_servicesScopePerRequest bool
	_toolsChangedDelegate    func()
	_promptsChangedDelegate  func()
	_resourceChangedDelegate func()
	_serverOnlyEndpointName  string
	EndpointName             string
	_started                 int32
	ServerCapabilities       *ServerCapabilities
	ClientCapabilities       *protocol.ClientCapabilities
	ClientInfo               *protocol.Implementation
	ServerOptions            McpServerOptions
}

func NewMcpServer(itransport protocol.ITransport, options McpServerOptions) *McpServer {
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
		_serverOnlyEndpointName:  fmt.Sprintf("Server (%s %s)", serverName, version),
		ServerOptions:            options,
		ClientInfo:               options.KnownClientInfo,
	}
	s.updateEndpointNameWithClientInfo()

	s.setInitializeHandler(&options)
	s.setToolsHandler(&options)
	s.setPromptsHandler(&options)
	s.setResourcesHandler(&options)
	s.setCompletionHandler(&options)
	s.setPingHandler()

	if options.Capabilities != nil && len(options.Capabilities.NotificationHandlers) > 0 {
		s.GetNotificationHandlers().RegisterRange(options.Capabilities.NotificationHandlers)
	}

	if t, ok := itransport.(*StreamableHttpServerTransport); !ok || !t.Stateless {
		if options.Capabilities != nil && options.Capabilities.Tools != nil && options.Capabilities.Tools.ToolCollection.Count() > 0 {
			s._toolsChangedDelegate = func() {
				s.SendMessage(context.Background(), protocol.NewJsonRpcNotification(protocol.NotificationMethods_ToolListChangedNotification, nil))
			}
			options.Capabilities.Tools.ToolCollection.OnChanged(s._toolsChangedDelegate)
		}

		if options.Capabilities != nil && options.Capabilities.Prompts != nil && options.Capabilities.Prompts.PromptCollection.Count() > 0 {
			s._promptsChangedDelegate = func() {
				s.SendMessage(context.Background(), protocol.NewJsonRpcNotification(protocol.NotificationMethods_PromptListChangedNotification, nil))
			}
			options.Capabilities.Prompts.PromptCollection.OnChanged(s._promptsChangedDelegate)
		}

		if options.Capabilities != nil && options.Capabilities.Resources != nil && options.Capabilities.Resources.ResourceCollection.Count() > 0 {
			s._resourceChangedDelegate = func() {
				s.SendMessage(context.Background(), protocol.NewJsonRpcNotification(protocol.NotificationMethods_ResourceListChangedNotification, nil))
			}
			options.Capabilities.Resources.ResourceCollection.OnChanged(s._resourceChangedDelegate)
		}

	}

	s.InitializeSession(itransport, true)
	return s
}
func (m *McpServer) updateEndpointNameWithClientInfo() {
	if m.ClientInfo == nil {
		return
	}

	m.EndpointName = fmt.Sprintf("%s, Client (%s %s)", m._serverOnlyEndpointName, m.ClientInfo.Name, m.ClientInfo.Version)
}

func (m *McpServer) setInitializeHandler(options *McpServerOptions) {
	shared.GenericRequestHandlerAdd(
		m.GetRequestHandlers(),
		protocol.RequestMethods_Initialize,
		func(ctx context.Context, request *protocol.InitializeRequestParams, tran protocol.ITransport) (*InitializeResult, error) {
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
	var listToolsHandler func(ctx context.Context, req RequestContext[*protocol.ListToolsRequestParams]) (*protocol.ListToolsResult, error)
	var callToolHandler func(ctx context.Context, req RequestContext[*protocol.CallToolRequestParams]) (*protocol.CallToolResult, error)
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
		listToolsHandler = func(ctx context.Context, req RequestContext[*protocol.ListToolsRequestParams]) (*protocol.ListToolsResult, error) {
			result := &protocol.ListToolsResult{
				Tools: make([]protocol.Tool, 0, tools.Count()),
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
		callToolHandler = func(ctx context.Context, req RequestContext[*protocol.CallToolRequestParams]) (*protocol.CallToolResult, error) {
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
			NotificationHandlers: map[string]protocol.NotificationHandler{},
		}
	} else {
		m.ServerCapabilities = options.Capabilities
	}

	setHandler(m, protocol.RequestMethods_ToolsList, listToolsHandler, nil, nil)
	setHandler(m, protocol.RequestMethods_ToolsCall, callToolHandler, nil, nil)
}

func (m *McpServer) setPromptsHandler(options *McpServerOptions) {
	var promptsCapability *PromptsCapability
	var listPromptsHandler func(context.Context, RequestContext[*protocol.ListPromptsRequestParams]) (*protocol.ListPromptsResult, error)
	var getPromptHandler func(context.Context, RequestContext[*protocol.GetPromptRequestParams]) (*protocol.GetPromptResult, error)
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
		listPromptsHandler = func(ctx context.Context, request RequestContext[*protocol.ListPromptsRequestParams]) (*protocol.ListPromptsResult, error) {
			result := &protocol.ListPromptsResult{
				Prompts: make([]protocol.Prompt, 0, prompts.Count()),
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
		getPromptHandler = func(ctx context.Context, request RequestContext[*protocol.GetPromptRequestParams]) (*protocol.GetPromptResult, error) {
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
			NotificationHandlers: map[string]protocol.NotificationHandler{},
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

	setHandler(m, protocol.RequestMethods_PromptsList, listPromptsHandler, nil, nil)
	setHandler(m, protocol.RequestMethods_PromptsGet, getPromptHandler, nil, nil)
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
		listResourcesHandler = func(context.Context, RequestContext[*protocol.ListResourcesRequestParams]) (*protocol.ListResourcesResult, error) {
			return &protocol.ListResourcesResult{}, nil
		}
	}
	if listResourceTemplatesHandler == nil {
		listResourceTemplatesHandler = func(context.Context, RequestContext[*protocol.ListResourceTemplatesRequestParams]) (*protocol.ListResourceTemplatesResult, error) {
			return &protocol.ListResourceTemplatesResult{}, nil
		}
	}

	setHandler(m, protocol.RequestMethods_ResourcesList, listResourcesHandler, nil, nil)
	setHandler(m, protocol.RequestMethods_ResourcesRead, readResourceHandler, nil, nil)
	setHandler(m, protocol.RequestMethods_ResourcesTemplatesList, listResourceTemplatesHandler, nil, nil)

	if resourcesCapability.Subscribe == nil && !*resourcesCapability.Subscribe {
		return
	}

	subscribeHandler := resourcesCapability.SubscribeToResourcesHandler
	unsubscribeHandler := resourcesCapability.UnsubscribeFromResourcesHandler
	if subscribeHandler == nil || unsubscribeHandler == nil {
		return
	}

	setHandler(m, protocol.RequestMethods_ResourcesSubscribe, subscribeHandler, nil, nil)
	setHandler(m, protocol.RequestMethods_ResourcesUnsubscribe, unsubscribeHandler, nil, nil)
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

	setHandler(m, protocol.RequestMethods_CompletionComplete, completeHandler, nil, nil)
}

func (m *McpServer) setPingHandler() {
	setHandler(
		m,
		protocol.RequestMethods_Ping,
		func(context.Context, RequestContext[*protocol.PingRequest]) (*protocol.PingResult, error) {
			return &protocol.PingResult{}, nil
		}, nil, nil)
}

// GetClientCapabilities implements IMcpServer.
func (m *McpServer) GetClientCapabilities() *protocol.ClientCapabilities {
	return m.ClientCapabilities
}

// GetClientInfo implements IMcpServer.
func (m *McpServer) GetClientInfo() *protocol.Implementation {
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

func (e *McpServer) RequestRoots(ctx context.Context, request protocol.ListRootsRequestParams) (*protocol.ListRootsResult, error) {
	if err := throwIfRootsUnsupported(e); err != nil {
		return nil, err
	}

	if e.GetClientCapabilities() == nil || e.GetClientCapabilities().Roots == nil {
		return nil, fmt.Errorf("client capabilities roots not set")
	}
	req := protocol.NewJsonRpcRequest(protocol.RequestMethods_RootsList, request, nil)
	resp, err := e.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	var result protocol.ListRootsResult
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
		func(ctx context.Context, request *TRequest, destinationTransport protocol.ITransport) (*TResponse, error) {
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
	destinationTransport protocol.ITransport,
) (*TResult, error) {
	// TODO: handle _servicesScopePerRequest
	svr := NewDestinationBoundMcpServer(m, destinationTransport)
	return handler(ctx, RequestContext[TParams]{Params: args, Server: svr})
}

func (e *McpServer) Elicit(ctx context.Context, request protocol.ElicitRequestParams) (*protocol.ElicitResult, error) {
	if err := throwIfElicitationUnsupported(e); err != nil {
		return nil, err
	}

	req := protocol.NewJsonRpcRequest(protocol.RequestMethods_ElicitationCreate, request, nil)
	resp, err := e.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	var result protocol.ElicitResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (e *McpServer) Sample(ctx context.Context, request protocol.CreateMessageRequestParams) (*protocol.CreateMessageResult, error) {
	if err := throwIfSamplingUnsupported(e); err != nil {
		return nil, err
	}

	req := protocol.NewJsonRpcRequest(protocol.RequestMethods_SamplingCreateMessage, request, nil)
	resp, err := e.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	var result protocol.CreateMessageResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (e *McpServer) SampleWithChatMessage(ctx context.Context, messages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	if err := throwIfSamplingUnsupported(e); err != nil {
		return nil, err
	}

	samplingMessages := []protocol.SamplingMessage{}
	var systemPrompt *strings.Builder
	for _, message := range messages {
		if message.Role == chatcompletion.RoleSystem {
			if systemPrompt == nil {
				systemPrompt = &strings.Builder{}
			} else {
				systemPrompt.WriteString("\n")
			}

			systemPrompt.WriteString(message.Text())
			continue
		}

		if message.Role == chatcompletion.RoleUser || message.Role == chatcompletion.RoleAssistant {
			role := protocol.RoleUser
			if message.Role == chatcompletion.RoleAssistant {
				role = protocol.RoleAssistant
			}
			for _, content := range message.Contents {
				switch con := content.(type) {
				case *contents.TextContent:
					samplingMessages = append(samplingMessages, protocol.SamplingMessage{
						Content: protocol.Content{
							Type: "text",
							Text: &con.Text,
						},
						Role: role,
					})
				case *contents.DataContent:
					if con.MediaTypeStartsWith("image") || con.MediaTypeStartsWith("audio") {
						t := "image"
						if con.MediaTypeStartsWith("audio") {
							t = "audio"
						}
						decoded := base64.URLEncoding.EncodeToString(con.Data)
						samplingMessages = append(samplingMessages, protocol.SamplingMessage{
							Content: protocol.Content{
								Type:     t,
								MimeType: &con.MediaType,
								Data:     &decoded,
							},
							Role: role,
						})
					}
				}
			}
		}
	}

	var modelPreferences protocol.ModelPreferences
	if options != nil && options.ModelId != nil {
		modelPreferences = protocol.ModelPreferences{
			Hints: []protocol.ModelHint{{
				Name: options.ModelId,
			}},
		}
	}

	systemPromptString := systemPrompt.String()
	request := protocol.CreateMessageRequestParams{
		RequestParams:    protocol.RequestParams{},
		MaxTokens:        options.MaxOutputTokens,
		Messages:         samplingMessages,
		Metadata:         nil,
		ModelPreferences: modelPreferences,
		StopSequences:    options.StopSequences,
		SystemPrompt:     &systemPromptString,
		Temperature:      options.Temperature,
	}

	result, err := e.Sample(ctx, request)
	if err != nil {
		return nil, err
	}

	message := &chatcompletion.ChatMessage{
		Contents: []contents.IAIContent{mcp.ContentToAIContent(result.Content)},
	}
	if result.Role == protocol.RoleUser {
		message.Role = chatcompletion.RoleUser
	}
	if result.Role == protocol.RoleAssistant {
		message.Role = chatcompletion.RoleAssistant
	}
	resp := chatcompletion.NewChatResponse(nil, message)
	if result.StopReason != nil {
		if *result.StopReason == "maxTokens" {
			t := chatcompletion.ReasonLength
			resp.FinishReason = &t
		} else {

			t := chatcompletion.ReasonStop
			resp.FinishReason = &t
		}
	}
	return resp, nil
}

func throwIfSamplingUnsupported(server IMcpServer) error {
	if server.GetClientCapabilities() == nil || server.GetClientCapabilities().Sampling == nil {
		if server.GetMcpServerOptions() != nil && server.GetMcpServerOptions().KnownClientInfo != nil {
			return fmt.Errorf("sampling is not supported in stateless mode")
		}

		return fmt.Errorf("client does not support sampling")
	}
	return nil
}

func throwIfRootsUnsupported(server IMcpServer) error {
	if server.GetClientCapabilities() == nil || server.GetClientCapabilities().Roots == nil {
		if server.GetMcpServerOptions() != nil && server.GetMcpServerOptions().KnownClientInfo != nil {
			return fmt.Errorf("roots are not supported in stateless mode")
		}

		return fmt.Errorf("client does not support roots")
	}
	return nil
}

func throwIfElicitationUnsupported(server IMcpServer) error {
	if server.GetClientCapabilities() == nil || server.GetClientCapabilities().Elicitation == nil {
		if server.GetMcpServerOptions() != nil && server.GetMcpServerOptions().KnownClientInfo != nil {
			return fmt.Errorf("elicitation is not supported in stateless mode")
		}

		return fmt.Errorf("client does not support elicitation requests")
	}
	return nil
}
