package client

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"

	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/mcp/protocol/messages"
	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/mcp/protocol/types"
	"github.com/futugyou/yomawari/mcp/server"
	"github.com/futugyou/yomawari/mcp/shared"
	"github.com/google/uuid"
)

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
				messages.RequestMethods_SamplingCreateMessage,
				func(ctx context.Context, request *types.CreateMessageRequestParams) (*types.CreateMessageResult, error) {

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
				messages.RequestMethods_RootsList,
				func(ctx context.Context, request *types.ListRootsRequestParams) (*types.ListRootsResult, error) {
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

	jsonRpcRequest := messages.NewJsonRpcRequest(messages.RequestMethods_Initialize, params, nil)
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
	return m.SendMessage(ctx, messages.NewJsonRpcNotification(messages.NotificationMethods_InitializedNotification, nil))
}

// CallTool implements IMcpClient.
func (m *McpClient) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}, reporter shared.IProgressReporter) (*types.CallToolResult, error) {
	params := types.CallToolRequestParams{
		RequestParams: types.RequestParams{},
		Name:          toolName,
		Arguments:     arguments,
	}

	if reporter != nil {
		progressToken := messages.NewProgressTokenFromString(uuid.New().String())
		var handler shared.NotificationHandler = func(ctx context.Context, notification *messages.JsonRpcNotification) error {
			var pn messages.ProgressNotification
			if err := json.Unmarshal(notification.Params, &pn); err != nil {
				return err
			}
			if pn.ProgressToken != nil && *pn.ProgressToken == progressToken {
				reporter.Report(*pn.Progress)
			}
			return nil
		}
		m.RegisterNotificationHandler(messages.NotificationMethods_ProgressNotification, handler)
		params.Meta = &types.RequestParamsMetadata{ProgressToken: &progressToken}
	}

	jsonRpcRequest := messages.NewJsonRpcRequest(messages.RequestMethods_Initialize, params, nil)
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
	panic("unimplemented")
}

// EnumeratePrompts implements IMcpClient.
func (m *McpClient) EnumeratePrompts(ctx context.Context, client IMcpClient) (<-chan McpClientPrompt, <-chan error) {
	panic("unimplemented")
}

// EnumerateResourceTemplates implements IMcpClient.
func (m *McpClient) EnumerateResourceTemplates(ctx context.Context, client IMcpClient) (<-chan types.ResourceTemplate, <-chan error) {
	panic("unimplemented")
}

// EnumerateResources implements IMcpClient.
func (m *McpClient) EnumerateResources(ctx context.Context, client IMcpClient) (<-chan types.Resource, <-chan error) {
	panic("unimplemented")
}

// EnumerateTools implements IMcpClient.
func (m *McpClient) EnumerateTools(ctx context.Context) (<-chan McpClientTool, <-chan error) {
	panic("unimplemented")
}

// GetPrompt implements IMcpClient.
func (m *McpClient) GetPrompt(ctx context.Context, name string, arguments map[string]interface{}) (*types.GetPromptResult, error) {
	panic("unimplemented")
}

// GetServerCapabilities implements IMcpClient.
func (m *McpClient) GetServerCapabilities() *server.ServerCapabilities {
	panic("unimplemented")
}

// GetServerInfo implements IMcpClient.
func (m *McpClient) GetServerInfo() *types.Implementation {
	panic("unimplemented")
}

// GetServerInstructions implements IMcpClient.
func (m *McpClient) GetServerInstructions() *string {
	panic("unimplemented")
}

// ListPrompts implements IMcpClient.
func (m *McpClient) ListPrompts(ctx context.Context, client IMcpClient) ([]McpClientPrompt, error) {
	panic("unimplemented")
}

// ListResourceTemplates implements IMcpClient.
func (m *McpClient) ListResourceTemplates(ctx context.Context, client IMcpClient) ([]types.ResourceTemplate, error) {
	panic("unimplemented")
}

// ListResources implements IMcpClient.
func (m *McpClient) ListResources(ctx context.Context, client IMcpClient) ([]types.Resource, error) {
	panic("unimplemented")
}

// ListTools implements IMcpClient.
func (m *McpClient) ListTools(ctx context.Context) ([]McpClientTool, error) {
	panic("unimplemented")
}

// Ping implements IMcpClient.
func (m *McpClient) Ping(ctx context.Context) error {
	panic("unimplemented")
}

// ReadResource implements IMcpClient.
func (m *McpClient) ReadResource(ctx context.Context, uri string) (*types.ReadResourceResult, error) {
	panic("unimplemented")
}

// ReadResourceWithUri implements IMcpClient.
func (m *McpClient) ReadResourceWithUri(ctx context.Context, uri url.URL) (*types.ReadResourceResult, error) {
	panic("unimplemented")
}

// SetLoggingLevel implements IMcpClient.
func (m *McpClient) SetLoggingLevel(ctx context.Context, level types.LoggingLevel) error {
	panic("unimplemented")
}

// SetLoggingLevelWithLogLevel implements IMcpClient.
func (m *McpClient) SetLoggingLevelWithLogLevel(ctx context.Context, level logger.LogLevel) error {
	panic("unimplemented")
}

// SubscribeToResource implements IMcpClient.
func (m *McpClient) SubscribeToResource(ctx context.Context, uri string) error {
	panic("unimplemented")
}

// SubscribeToResourceWithUri implements IMcpClient.
func (m *McpClient) SubscribeToResourceWithUri(ctx context.Context, uri url.URL) error {
	panic("unimplemented")
}

// UnsubscribeFromResource implements IMcpClient.
func (m *McpClient) UnsubscribeFromResource(ctx context.Context, uri string) error {
	panic("unimplemented")
}

// UnsubscribeFromResourceWithUri implements IMcpClient.
func (m *McpClient) UnsubscribeFromResourceWithUri(ctx context.Context, uri url.URL) error {
	panic("unimplemented")
}
