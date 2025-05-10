package client

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"

	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/mcp/protocol/types"
	"github.com/futugyou/yomawari/mcp/server"
	"github.com/futugyou/yomawari/mcp/shared"
	"github.com/google/uuid"
)

var McpClientDefaultImplementation types.Implementation = types.Implementation{
	Name:    "McpClient",
	Version: "1.0.0",
}

var _ IMcpClient = (*McpClient)(nil)

type McpClient struct {
	*shared.BaseMcpEndpoint
	clientTransport  transport.IClientTransport
	options          McpClientOptions
	sessionTransport transport.ITransport
	reqHandlers      *shared.RequestHandlers
	notifHandlers    *shared.NotificationHandlers

	ctx        context.Context
	connectCts context.CancelFunc
	mu         sync.Mutex
	disposed   bool

	ServerCapabilities server.ServerCapabilities
	ServerInfo         types.Implementation
	ServerInstructions *string
	EndpointName       string
}

func CreateMcpClient(ctx context.Context, clientTransport transport.IClientTransport, options McpClientOptions) (*McpClient, error) {
	client := NewMcpClient(clientTransport, options)
	err := client.Connect(ctx)
	if err != nil {
		client.Dispose(ctx)
		return nil, err
	}
	return client, nil
}

func NewMcpClient(clientTransport transport.IClientTransport, options McpClientOptions) *McpClient {
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
				transport.RequestMethods_SamplingCreateMessage,
				func(ctx context.Context, request *types.CreateMessageRequestParams, tran transport.ITransport) (*types.CreateMessageResult, error) {
					var progres shared.IProgressReporter = &shared.NullProgress{}
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
				transport.RequestMethods_RootsList,
				func(ctx context.Context, request *types.ListRootsRequestParams, tran transport.ITransport) (*types.ListRootsResult, error) {
					return capabilities.Roots.RootsHandler(ctx, request)
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

	params := shared.InitializeRequestParams{
		ProtocolVersion: m.options.ProtocolVersion,
		Capabilities:    m.options.Capabilities,
	}

	if m.options.ClientInfo != nil {
		params.ClientInfo = *m.options.ClientInfo
	}

	jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_Initialize, params, nil)
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
	return m.SendMessage(ctx, transport.NewJsonRpcNotification(transport.NotificationMethods_InitializedNotification, nil))
}

// CallTool implements IMcpClient.
func (m *McpClient) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}, reporter shared.IProgressReporter) (*types.CallToolResult, error) {
	params := types.CallToolRequestParams{
		RequestParams: types.RequestParams{},
		Name:          toolName,
		Arguments:     arguments,
	}

	if reporter != nil {
		progressToken := transport.NewProgressTokenFromString(uuid.New().String())
		var handler shared.NotificationHandler = func(ctx context.Context, notification *transport.JsonRpcNotification) error {
			var pn transport.ProgressNotification
			if err := json.Unmarshal(notification.Params, &pn); err != nil {
				return err
			}
			if pn.ProgressToken != nil && *pn.ProgressToken == progressToken {
				reporter.Report(*pn.Progress)
			}
			return nil
		}
		m.RegisterNotificationHandler(transport.NotificationMethods_ProgressNotification, handler)
		params.Meta = &types.RequestParamsMetadata{ProgressToken: &progressToken}
	}

	jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_ToolsCall, params, nil)
	resp, err := m.SendRequest(ctx, jsonRpcRequest)
	if err != nil {
		return nil, err
	}
	var rsult types.CallToolResult
	if err := json.Unmarshal(resp.Result, &rsult); err != nil {
		return nil, err
	}
	return &rsult, nil
}

// Complete implements IMcpClient.
func (m *McpClient) Complete(ctx context.Context, reference types.Reference, argumentName string, argumentValue string) (*types.CompleteResult, error) {
	params := types.CompleteRequestParams{
		Ref: reference,
		Argument: types.Argument{
			Name:  argumentName,
			Value: argumentValue,
		},
	}
	jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_CompletionComplete, params, nil)
	resp, err := m.SendRequest(ctx, jsonRpcRequest)
	if err != nil {
		return nil, err
	}
	var rsult types.CompleteResult
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
			params := types.ListPromptsRequestParams{
				PaginatedRequestParams: types.PaginatedRequestParams{
					RequestParams: types.RequestParams{},
					Cursor:        cursor,
				},
			}
			jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_PromptsList, params, nil)
			resp, err := m.SendRequest(ctx, jsonRpcRequest)
			if err != nil {
				errCh <- err
				return
			}

			var promptResults types.ListPromptsResult
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
			params := types.ListResourceTemplatesRequestParams{
				PaginatedRequestParams: types.PaginatedRequestParams{
					RequestParams: types.RequestParams{},
					Cursor:        cursor,
				},
			}
			jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_ResourcesTemplatesList, params, nil)
			resp, err := m.SendRequest(ctx, jsonRpcRequest)
			if err != nil {
				errCh <- err
				return
			}

			var promptResults types.ListResourceTemplatesResult
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
			params := types.ListResourcesRequestParams{
				PaginatedRequestParams: types.PaginatedRequestParams{
					RequestParams: types.RequestParams{},
					Cursor:        cursor,
				},
			}
			jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_ResourcesList, params, nil)
			resp, err := m.SendRequest(ctx, jsonRpcRequest)
			if err != nil {
				errCh <- err
				return
			}

			var promptResults types.ListResourcesResult
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
			params := types.ListToolsRequestParams{
				PaginatedRequestParams: types.PaginatedRequestParams{
					RequestParams: types.RequestParams{},
					Cursor:        cursor,
				},
			}
			jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_ToolsList, params, nil)
			resp, err := m.SendRequest(ctx, jsonRpcRequest)
			if err != nil {
				errCh <- err
				return
			}

			var promptResults types.ListToolsResult
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
func (m *McpClient) GetPrompt(ctx context.Context, name string, arguments map[string]interface{}) (*types.GetPromptResult, error) {
	params := types.GetPromptRequestParams{
		RequestParams: types.RequestParams{},
		Name:          name,
		Arguments:     arguments,
	}
	jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_PromptsGet, params, nil)
	resp, err := m.SendRequest(ctx, jsonRpcRequest)
	if err != nil {
		return nil, err
	}
	var rsult types.GetPromptResult
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
func (m *McpClient) GetServerInfo() *types.Implementation {
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
		params := types.ListPromptsRequestParams{
			PaginatedRequestParams: types.PaginatedRequestParams{
				RequestParams: types.RequestParams{},
				Cursor:        cursor,
			},
		}
		jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_PromptsList, params, nil)
		resp, err := m.SendRequest(ctx, jsonRpcRequest)
		if err != nil {
			return nil, err
		}

		var promptResults types.ListPromptsResult
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
		params := types.ListResourceTemplatesRequestParams{
			PaginatedRequestParams: types.PaginatedRequestParams{
				RequestParams: types.RequestParams{},
				Cursor:        cursor,
			},
		}
		jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_ResourcesTemplatesList, params, nil)
		resp, err := m.SendRequest(ctx, jsonRpcRequest)
		if err != nil {
			return nil, err
		}

		var promptResults types.ListResourceTemplatesResult
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
		params := types.ListResourcesRequestParams{
			PaginatedRequestParams: types.PaginatedRequestParams{
				RequestParams: types.RequestParams{},
				Cursor:        cursor,
			},
		}
		jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_ResourcesList, params, nil)
		resp, err := m.SendRequest(ctx, jsonRpcRequest)
		if err != nil {
			return nil, err
		}

		var promptResults types.ListResourcesResult
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
		params := types.ListToolsRequestParams{
			PaginatedRequestParams: types.PaginatedRequestParams{
				RequestParams: types.RequestParams{},
				Cursor:        cursor,
			},
		}
		jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_ToolsList, params, nil)
		resp, err := m.SendRequest(ctx, jsonRpcRequest)
		if err != nil {
			return nil, err
		}

		var promptResults types.ListToolsResult
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
	params := types.PingRequest{}
	jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_Ping, params, nil)
	_, err := m.SendRequest(ctx, jsonRpcRequest)
	return err
}

// ReadResource implements IMcpClient.
func (m *McpClient) ReadResource(ctx context.Context, uri string) (*types.ReadResourceResult, error) {
	params := types.ReadResourceRequestParams{
		RequestParams: types.RequestParams{},
		Uri:           uri,
	}
	jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_ResourcesRead, params, nil)
	resp, err := m.SendRequest(ctx, jsonRpcRequest)
	if err != nil {
		return nil, err
	}
	var rsult types.ReadResourceResult
	if err := json.Unmarshal(resp.Result, &rsult); err != nil {
		return nil, err
	}
	return &rsult, nil
}

// ReadResourceWithUri implements IMcpClient.
func (m *McpClient) ReadResourceWithUri(ctx context.Context, uri url.URL) (*types.ReadResourceResult, error) {
	return m.ReadResource(ctx, uri.String())
}

// ReadResourceWithUriAndArguments implements IMcpClient.
func (m *McpClient) ReadResourceWithUriAndArguments(ctx context.Context, uriTemplate string, arguments map[string]interface{}) (*types.ReadResourceResult, error) {
	url, err := shared.FormatUri(uriTemplate, arguments)
	if err != nil {
		return nil, err
	}

	return m.ReadResource(ctx, url)
}

// SetLoggingLevel implements IMcpClient.
func (m *McpClient) SetLoggingLevel(ctx context.Context, level types.LoggingLevel) error {
	params := types.SetLevelRequestParams{
		RequestParams: types.RequestParams{},
		Level:         level,
	}
	jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_LoggingSetLevel, params, nil)
	_, err := m.SendRequest(ctx, jsonRpcRequest)
	return err
}

// SetLoggingLevelWithLogLevel implements IMcpClient.
func (m *McpClient) SetLoggingLevelWithLogLevel(ctx context.Context, level logger.LogLevel) error {
	return m.SetLoggingLevel(ctx, types.LoggingLevel(level))
}

// SubscribeToResource implements IMcpClient.
func (m *McpClient) SubscribeToResource(ctx context.Context, uri string) error {
	params := types.SubscribeRequestParams{
		RequestParams: types.RequestParams{},
		Uri:           &uri,
	}
	jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_ResourcesSubscribe, params, nil)
	_, err := m.SendRequest(ctx, jsonRpcRequest)
	return err
}

// SubscribeToResourceWithUri implements IMcpClient.
func (m *McpClient) SubscribeToResourceWithUri(ctx context.Context, uri url.URL) error {
	return m.SubscribeToResource(ctx, uri.String())
}

// UnsubscribeFromResource implements IMcpClient.
func (m *McpClient) UnsubscribeFromResource(ctx context.Context, uri string) error {
	params := types.UnsubscribeRequestParams{
		RequestParams: types.RequestParams{},
		Uri:           &uri,
	}
	jsonRpcRequest := transport.NewJsonRpcRequest(transport.RequestMethods_ResourcesUnsubscribe, params, nil)
	_, err := m.SendRequest(ctx, jsonRpcRequest)
	return err
}

// UnsubscribeFromResourceWithUri implements IMcpClient.
func (m *McpClient) UnsubscribeFromResourceWithUri(ctx context.Context, uri url.URL) error {
	return m.UnsubscribeFromResource(ctx, uri.String())
}
