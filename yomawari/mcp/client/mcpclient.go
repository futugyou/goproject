package client

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"

	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/mcp/protocol"
	"github.com/futugyou/yomawari/mcp/server"
	"github.com/futugyou/yomawari/mcp/shared"
	"github.com/google/uuid"
)

var McpClientDefaultImplementation protocol.Implementation = protocol.Implementation{
	Name:    "McpClient",
	Version: "1.0.0",
}

var _ IMcpClient = (*McpClient)(nil)

type McpClient struct {
	*shared.BaseMcpEndpoint
	clientTransport  IClientTransport
	options          McpClientOptions
	sessionTransport protocol.ITransport
	reqHandlers      *shared.RequestHandlers
	notifHandlers    *shared.NotificationHandlers

	ctx        context.Context
	connectCts context.CancelFunc
	mu         sync.Mutex
	disposed   bool

	ServerCapabilities server.ServerCapabilities
	ServerInfo         protocol.Implementation
	ServerInstructions *string
	EndpointName       string
}

func CreateMcpClient(ctx context.Context, clientTransport IClientTransport, options McpClientOptions) (*McpClient, error) {
	client := NewMcpClient(clientTransport, options)
	err := client.Connect(ctx)
	if err != nil {
		client.Dispose(ctx)
		return nil, err
	}
	return client, nil
}

func NewMcpClient(clientTransport IClientTransport, options McpClientOptions) *McpClient {
	client := &McpClient{
		BaseMcpEndpoint: shared.NewBaseMcpEndpoint(),
		clientTransport: clientTransport,
		options:         options,
		EndpointName:    clientTransport.GetName(),
		reqHandlers:     shared.NewRequestHandlers(),
		notifHandlers:   shared.NewNotificationHandlers(),
	}

	capabilities := options.Capabilities
	if capabilities != nil {
		notificationHandlers := capabilities.NotificationHandlers
		if notificationHandlers != nil {
			client.notifHandlers.RegisterRange(notificationHandlers)
		}
		samplingCapability := capabilities.Sampling
		if samplingCapability != nil && samplingCapability.SamplingHandler != nil {
			shared.GenericRequestHandlerAdd(
				client.reqHandlers,
				protocol.RequestMethods_SamplingCreateMessage,
				func(ctx context.Context, request *protocol.CreateMessageRequestParams, tran protocol.ITransport) (*protocol.CreateMessageResult, error) {
					var progres protocol.IProgressReporter = &shared.NullProgress{}
					if request.Meta != nil && request.Meta.ProgressToken != nil {
						progres = shared.NewTokenProgress(client, *request.Meta.ProgressToken)
					}
					return samplingCapability.SamplingHandler(ctx, request, progres)
				},
				nil,
				nil,
			)
		}

		if capabilities.Roots != nil && capabilities.Roots.RootsHandler != nil {
			shared.GenericRequestHandlerAdd(
				client.reqHandlers,
				protocol.RequestMethods_RootsList,
				func(ctx context.Context, request *protocol.ListRootsRequestParams, tran protocol.ITransport) (*protocol.ListRootsResult, error) {
					return capabilities.Roots.RootsHandler(ctx, request)
				},
				nil,
				nil,
			)
		}

		if capabilities.Elicitation != nil && capabilities.Elicitation.ElicitationHandler != nil {
			shared.GenericRequestHandlerAdd(
				client.reqHandlers,
				protocol.RequestMethods_ElicitationCreate,
				func(ctx context.Context, request *protocol.ElicitRequestParams, tran protocol.ITransport) (*protocol.ElicitResult, error) {
					return capabilities.Elicitation.ElicitationHandler(ctx, request)
				},
				nil,
				nil,
			)
		}
	}
	return client
}

func (e *McpClient) Dispose(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.disposed {
		return nil
	}
	e.disposed = true

	if e.connectCts != nil {
		e.connectCts()
	}

	defer e.sessionTransport.Close()

	return e.BaseMcpEndpoint.Dispose(ctx)
}

func (m *McpClient) Connect(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	m.ctx = ctx
	m.connectCts = cancel
	sessionTransport, err := m.clientTransport.Connect(ctx)
	if err != nil {
		return err
	}
	m.sessionTransport = sessionTransport

	m.InitializeSession(sessionTransport, false)
	// We don't want the ConnectAsync token to cancel the session after we've successfully connected.
	// The base class handles cleaning up the session in DisposeAsync without our help.
	m.StartSession(context.Background(), sessionTransport)
	ctx, cancel = context.WithTimeout(ctx, m.options.InitializationTimeout)
	defer cancel()

	params := protocol.InitializeRequestParams{
		ProtocolVersion: m.options.ProtocolVersion,
		Capabilities:    m.options.Capabilities,
	}

	if m.options.ClientInfo != nil {
		params.ClientInfo = *m.options.ClientInfo
	}

	jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_Initialize, params, nil)
	initializeResponse, err := m.SendRequest(ctx, jsonRpcRequest)
	if err != nil {
		return err
	}
	var initializeResult server.InitializeResult
	if err := json.Unmarshal(initializeResponse.Result, &initializeResult); err != nil {
		return err
	}

	m.ServerCapabilities = initializeResult.Capabilities
	m.ServerInfo = initializeResult.ServerInfo
	m.ServerInstructions = &initializeResult.Instructions
	return m.SendMessage(ctx, protocol.NewJsonRpcNotification(protocol.NotificationMethods_InitializedNotification, nil))
}

// CallTool implements IMcpClient.
func (m *McpClient) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}, reporter protocol.IProgressReporter) (*protocol.CallToolResult, error) {
	params := protocol.CallToolRequestParams{
		RequestParams: protocol.RequestParams{},
		Name:          toolName,
		Arguments:     arguments,
	}

	if reporter != nil {
		progressToken := protocol.NewProgressTokenFromString(uuid.New().String())
		var handler protocol.NotificationHandler = func(ctx context.Context, notification *protocol.JsonRpcNotification) error {
			var pn protocol.ProgressNotification
			if err := json.Unmarshal(notification.Params, &pn); err != nil {
				return err
			}
			if pn.ProgressToken != nil && *pn.ProgressToken == progressToken {
				reporter.Report(*pn.Progress)
			}
			return nil
		}
		m.RegisterNotificationHandler(protocol.NotificationMethods_ProgressNotification, handler)
		params.Meta = &protocol.RequestParamsMetadata{ProgressToken: &progressToken}
	}

	jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_ToolsCall, params, nil)
	resp, err := m.SendRequest(ctx, jsonRpcRequest)
	if err != nil {
		return nil, err
	}
	var rsult protocol.CallToolResult
	if err := json.Unmarshal(resp.Result, &rsult); err != nil {
		return nil, err
	}
	return &rsult, nil
}

// Complete implements IMcpClient.
func (m *McpClient) Complete(ctx context.Context, reference protocol.Reference, argumentName string, argumentValue string) (*protocol.CompleteResult, error) {
	params := protocol.CompleteRequestParams{
		Ref: reference,
		Argument: protocol.Argument{
			Name:  argumentName,
			Value: argumentValue,
		},
	}
	jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_CompletionComplete, params, nil)
	resp, err := m.SendRequest(ctx, jsonRpcRequest)
	if err != nil {
		return nil, err
	}
	var rsult protocol.CompleteResult
	if err := json.Unmarshal(resp.Result, &rsult); err != nil {
		return nil, err
	}
	return &rsult, nil
}

// EnumeratePrompts implements IMcpClient.
func (m *McpClient) EnumeratePrompts(ctx context.Context, client IMcpClient) (<-chan McpClientPrompt, <-chan error) {
	promptsCh := make(chan McpClientPrompt)
	errCh := make(chan error, 1)

	go func() {
		defer close(promptsCh)
		defer close(errCh)

		var cursor *string
		for {
			params := protocol.ListPromptsRequestParams{
				PaginatedRequestParams: protocol.PaginatedRequestParams{
					RequestParams: protocol.RequestParams{},
					Cursor:        cursor,
				},
			}
			jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_PromptsList, params, nil)
			resp, err := m.SendRequest(ctx, jsonRpcRequest)
			if err != nil {
				errCh <- err
				return
			}

			var promptResults protocol.ListPromptsResult
			if err := json.Unmarshal(resp.Result, &promptResults); err != nil {
				errCh <- err
				return
			}

			for _, prompt := range promptResults.Prompts {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case promptsCh <- *NewMcpClientPrompt(prompt, m):
				}
			}

			if promptResults.NextCursor == nil {
				break
			}
			cursor = promptResults.NextCursor
		}
	}()

	return promptsCh, errCh
}

// EnumerateResourceTemplates implements IMcpClient.
func (m *McpClient) EnumerateResourceTemplates(ctx context.Context, client IMcpClient) (<-chan McpClientResourceTemplate, <-chan error) {
	promptsCh := make(chan McpClientResourceTemplate)
	errCh := make(chan error, 1)

	go func() {
		defer close(promptsCh)
		defer close(errCh)

		var cursor *string
		for {
			params := protocol.ListResourceTemplatesRequestParams{
				PaginatedRequestParams: protocol.PaginatedRequestParams{
					RequestParams: protocol.RequestParams{},
					Cursor:        cursor,
				},
			}
			jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_ResourcesTemplatesList, params, nil)
			resp, err := m.SendRequest(ctx, jsonRpcRequest)
			if err != nil {
				errCh <- err
				return
			}

			var promptResults protocol.ListResourceTemplatesResult
			if err := json.Unmarshal(resp.Result, &promptResults); err != nil {
				errCh <- err
				return
			}

			for _, prompt := range promptResults.ResourceTemplates {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					t := NewMcpClientResourceTemplate(m, prompt)
					promptsCh <- *t
				}
			}

			if promptResults.NextCursor == nil {
				break
			}
			cursor = promptResults.NextCursor
		}
	}()

	return promptsCh, errCh
}

// EnumerateResources implements IMcpClient.
func (m *McpClient) EnumerateResources(ctx context.Context, client IMcpClient) (<-chan McpClientResource, <-chan error) {
	promptsCh := make(chan McpClientResource)
	errCh := make(chan error, 1)

	go func() {
		defer close(promptsCh)
		defer close(errCh)

		var cursor *string
		for {
			params := protocol.ListResourcesRequestParams{
				PaginatedRequestParams: protocol.PaginatedRequestParams{
					RequestParams: protocol.RequestParams{},
					Cursor:        cursor,
				},
			}
			jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_ResourcesList, params, nil)
			resp, err := m.SendRequest(ctx, jsonRpcRequest)
			if err != nil {
				errCh <- err
				return
			}

			var promptResults protocol.ListResourcesResult
			if err := json.Unmarshal(resp.Result, &promptResults); err != nil {
				errCh <- err
				return
			}

			for _, prompt := range promptResults.Resources {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				default:
					t := NewMcpClientResource(m, prompt)
					promptsCh <- *t
				}
			}

			if promptResults.NextCursor == nil {
				break
			}
			cursor = promptResults.NextCursor
		}
	}()

	return promptsCh, errCh
}

// EnumerateTools implements IMcpClient.
func (m *McpClient) EnumerateTools(ctx context.Context) (<-chan McpClientTool, <-chan error) {
	promptsCh := make(chan McpClientTool)
	errCh := make(chan error, 1)

	go func() {
		defer close(promptsCh)
		defer close(errCh)

		var cursor *string
		for {
			params := protocol.ListToolsRequestParams{
				PaginatedRequestParams: protocol.PaginatedRequestParams{
					RequestParams: protocol.RequestParams{},
					Cursor:        cursor,
				},
			}
			jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_ToolsList, params, nil)
			resp, err := m.SendRequest(ctx, jsonRpcRequest)
			if err != nil {
				errCh <- err
				return
			}

			var promptResults protocol.ListToolsResult
			if err := json.Unmarshal(resp.Result, &promptResults); err != nil {
				errCh <- err
				return
			}

			for _, prompt := range promptResults.Tools {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case promptsCh <- *NewMcpClientTool(m, prompt.Name, prompt.Description, prompt):
				}
			}

			if promptResults.NextCursor == nil {
				break
			}
			cursor = promptResults.NextCursor
		}
	}()

	return promptsCh, errCh
}

// GetPrompt implements IMcpClient.
func (m *McpClient) GetPrompt(ctx context.Context, name string, arguments map[string]interface{}) (*protocol.GetPromptResult, error) {
	params := protocol.GetPromptRequestParams{
		RequestParams: protocol.RequestParams{},
		Name:          name,
		Arguments:     arguments,
	}
	jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_PromptsGet, params, nil)
	resp, err := m.SendRequest(ctx, jsonRpcRequest)
	if err != nil {
		return nil, err
	}
	var rsult protocol.GetPromptResult
	if err := json.Unmarshal(resp.Result, &rsult); err != nil {
		return nil, err
	}
	return &rsult, nil
}

// GetServerCapabilities implements IMcpClient.
func (m *McpClient) GetServerCapabilities() *server.ServerCapabilities {
	return &m.ServerCapabilities
}

// GetServerInfo implements IMcpClient.
func (m *McpClient) GetServerInfo() *protocol.Implementation {
	return &m.ServerInfo
}

// GetServerInstructions implements IMcpClient.
func (m *McpClient) GetServerInstructions() *string {
	return m.ServerInstructions
}

// ListPrompts implements IMcpClient.
func (m *McpClient) ListPrompts(ctx context.Context, client IMcpClient) ([]McpClientPrompt, error) {
	prompts := []McpClientPrompt{}
	var cursor *string
	for {
		params := protocol.ListPromptsRequestParams{
			PaginatedRequestParams: protocol.PaginatedRequestParams{
				RequestParams: protocol.RequestParams{},
				Cursor:        cursor,
			},
		}
		jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_PromptsList, params, nil)
		resp, err := m.SendRequest(ctx, jsonRpcRequest)
		if err != nil {
			return nil, err
		}

		var promptResults protocol.ListPromptsResult
		if err := json.Unmarshal(resp.Result, &promptResults); err != nil {
			return nil, err
		}

		for _, prompt := range promptResults.Prompts {
			prompts = append(prompts, *NewMcpClientPrompt(prompt, m))
		}

		if promptResults.NextCursor == nil {
			break
		}
		cursor = promptResults.NextCursor
	}
	return prompts, nil
}

// ListResourceTemplates implements IMcpClient.
func (m *McpClient) ListResourceTemplates(ctx context.Context, client IMcpClient) ([]McpClientResourceTemplate, error) {
	prompts := []McpClientResourceTemplate{}
	var cursor *string
	for {
		params := protocol.ListResourceTemplatesRequestParams{
			PaginatedRequestParams: protocol.PaginatedRequestParams{
				RequestParams: protocol.RequestParams{},
				Cursor:        cursor,
			},
		}
		jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_ResourcesTemplatesList, params, nil)
		resp, err := m.SendRequest(ctx, jsonRpcRequest)
		if err != nil {
			return nil, err
		}

		var promptResults protocol.ListResourceTemplatesResult
		if err := json.Unmarshal(resp.Result, &promptResults); err != nil {
			return nil, err
		}

		for _, v := range promptResults.ResourceTemplates {
			t := NewMcpClientResourceTemplate(m, v)
			prompts = append(prompts, *t)
		}

		if promptResults.NextCursor == nil {
			break
		}
		cursor = promptResults.NextCursor
	}
	return prompts, nil
}

// ListResources implements IMcpClient.
func (m *McpClient) ListResources(ctx context.Context, client IMcpClient) ([]McpClientResource, error) {
	prompts := []McpClientResource{}
	var cursor *string
	for {
		params := protocol.ListResourcesRequestParams{
			PaginatedRequestParams: protocol.PaginatedRequestParams{
				RequestParams: protocol.RequestParams{},
				Cursor:        cursor,
			},
		}
		jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_ResourcesList, params, nil)
		resp, err := m.SendRequest(ctx, jsonRpcRequest)
		if err != nil {
			return nil, err
		}

		var promptResults protocol.ListResourcesResult
		if err := json.Unmarshal(resp.Result, &promptResults); err != nil {
			return nil, err
		}

		for _, v := range promptResults.Resources {
			t := NewMcpClientResource(m, v)
			prompts = append(prompts, *t)
		}

		if promptResults.NextCursor == nil {
			break
		}
		cursor = promptResults.NextCursor
	}
	return prompts, nil
}

// ListTools implements IMcpClient.
func (m *McpClient) ListTools(ctx context.Context) ([]McpClientTool, error) {
	prompts := []McpClientTool{}
	var cursor *string
	for {
		params := protocol.ListToolsRequestParams{
			PaginatedRequestParams: protocol.PaginatedRequestParams{
				RequestParams: protocol.RequestParams{},
				Cursor:        cursor,
			},
		}
		jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_ToolsList, params, nil)
		resp, err := m.SendRequest(ctx, jsonRpcRequest)
		if err != nil {
			return nil, err
		}

		var promptResults protocol.ListToolsResult
		if err := json.Unmarshal(resp.Result, &promptResults); err != nil {
			return nil, err
		}

		for _, v := range promptResults.Tools {
			prompts = append(prompts, *NewMcpClientTool(m, v.Name, v.Description, v))
		}

		if promptResults.NextCursor == nil {
			break
		}
		cursor = promptResults.NextCursor
	}
	return prompts, nil
}

// Ping implements IMcpClient.
func (m *McpClient) Ping(ctx context.Context) error {
	params := protocol.PingRequest{}
	jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_Ping, params, nil)
	_, err := m.SendRequest(ctx, jsonRpcRequest)
	return err
}

// ReadResource implements IMcpClient.
func (m *McpClient) ReadResource(ctx context.Context, uri string) (*protocol.ReadResourceResult, error) {
	params := protocol.ReadResourceRequestParams{
		RequestParams: protocol.RequestParams{},
		Uri:           uri,
	}
	jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_ResourcesRead, params, nil)
	resp, err := m.SendRequest(ctx, jsonRpcRequest)
	if err != nil {
		return nil, err
	}
	var rsult protocol.ReadResourceResult
	if err := json.Unmarshal(resp.Result, &rsult); err != nil {
		return nil, err
	}
	return &rsult, nil
}

// ReadResourceWithUri implements IMcpClient.
func (m *McpClient) ReadResourceWithUri(ctx context.Context, uri url.URL) (*protocol.ReadResourceResult, error) {
	return m.ReadResource(ctx, uri.String())
}

// ReadResourceWithUriAndArguments implements IMcpClient.
func (m *McpClient) ReadResourceWithUriAndArguments(ctx context.Context, uriTemplate string, arguments map[string]interface{}) (*protocol.ReadResourceResult, error) {
	url, err := shared.FormatUri(uriTemplate, arguments)
	if err != nil {
		return nil, err
	}

	return m.ReadResource(ctx, url)
}

// SetLoggingLevel implements IMcpClient.
func (m *McpClient) SetLoggingLevel(ctx context.Context, level protocol.LoggingLevel) error {
	params := protocol.SetLevelRequestParams{
		RequestParams: protocol.RequestParams{},
		Level:         level,
	}
	jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_LoggingSetLevel, params, nil)
	_, err := m.SendRequest(ctx, jsonRpcRequest)
	return err
}

// SetLoggingLevelWithLogLevel implements IMcpClient.
func (m *McpClient) SetLoggingLevelWithLogLevel(ctx context.Context, level logger.LogLevel) error {
	return m.SetLoggingLevel(ctx, protocol.LoggingLevel(level))
}

// SubscribeToResource implements IMcpClient.
func (m *McpClient) SubscribeToResource(ctx context.Context, uri string) error {
	params := protocol.SubscribeRequestParams{
		RequestParams: protocol.RequestParams{},
		Uri:           &uri,
	}
	jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_ResourcesSubscribe, params, nil)
	_, err := m.SendRequest(ctx, jsonRpcRequest)
	return err
}

// SubscribeToResourceWithUri implements IMcpClient.
func (m *McpClient) SubscribeToResourceWithUri(ctx context.Context, uri url.URL) error {
	return m.SubscribeToResource(ctx, uri.String())
}

// UnsubscribeFromResource implements IMcpClient.
func (m *McpClient) UnsubscribeFromResource(ctx context.Context, uri string) error {
	params := protocol.UnsubscribeRequestParams{
		RequestParams: protocol.RequestParams{},
		Uri:           &uri,
	}
	jsonRpcRequest := protocol.NewJsonRpcRequest(protocol.RequestMethods_ResourcesUnsubscribe, params, nil)
	_, err := m.SendRequest(ctx, jsonRpcRequest)
	return err
}

// UnsubscribeFromResourceWithUri implements IMcpClient.
func (m *McpClient) UnsubscribeFromResourceWithUri(ctx context.Context, uri url.URL) error {
	return m.UnsubscribeFromResource(ctx, uri.String())
}
