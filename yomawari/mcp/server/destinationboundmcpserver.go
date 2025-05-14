package server

import (
	"context"
	"fmt"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/mcp/protocol/types"
	"github.com/futugyou/yomawari/mcp/shared"
)

var _ IMcpServer = (*DestinationBoundMcpServer)(nil)

type DestinationBoundMcpServer struct {
	server    *McpServer
	transport transport.ITransport
}

func NewDestinationBoundMcpServer(server *McpServer, transport transport.ITransport) *DestinationBoundMcpServer {
	return &DestinationBoundMcpServer{
		server:    server,
		transport: transport,
	}
}

// AsSamplingChatClient implements IMcpServer.
func (d *DestinationBoundMcpServer) AsSamplingChatClient() (chatcompletion.IChatClient, error) {
	if d.GetClientCapabilities() == nil || d.GetClientCapabilities().Sampling == nil {
		return nil, fmt.Errorf("client capabilities sampling not set")
	}
	return NewSamplingChatClient(d), nil
}

// Dispose implements IMcpServer.
func (d *DestinationBoundMcpServer) Dispose(ctx context.Context) error {
	return d.server.Dispose(ctx)
}

// GetClientCapabilities implements IMcpServer.
func (d *DestinationBoundMcpServer) GetClientCapabilities() *transport.ClientCapabilities {
	return d.server.GetClientCapabilities()
}

// GetClientInfo implements IMcpServer.
func (d *DestinationBoundMcpServer) GetClientInfo() *types.Implementation {
	return d.server.GetClientInfo()
}

// GetEndpointName implements IMcpServer.
func (d *DestinationBoundMcpServer) GetEndpointName() string {
	return d.server.EndpointName
}

// GetMcpServerOptions implements IMcpServer.
func (d *DestinationBoundMcpServer) GetMcpServerOptions() *McpServerOptions {
	return d.server.GetMcpServerOptions()
}

// GetMessageProcessingTask implements IMcpServer.
func (d *DestinationBoundMcpServer) GetMessageProcessingTask() <-chan struct{} {
	return d.server.GetMessageProcessingTask()
}

// NotifyProgress implements IMcpServer.
func (d *DestinationBoundMcpServer) NotifyProgress(ctx context.Context, progressToken types.ProgressToken, progress transport.ProgressNotificationValue) error {
	return d.server.NotifyProgress(ctx, progressToken, progress)
}

// RegisterNotificationHandler implements IMcpServer.
func (d *DestinationBoundMcpServer) RegisterNotificationHandler(method string, handler transport.NotificationHandler) *shared.RegistrationHandle {
	return d.server.RegisterNotificationHandler(method, handler)
}

// RequestRoots implements IMcpServer.
func (d *DestinationBoundMcpServer) RequestRoots(ctx context.Context, request types.ListRootsRequestParams) (*types.ListRootsResult, error) {
	return d.server.RequestRoots(ctx, request)
}

// RequestSampling implements IMcpServer.
func (d *DestinationBoundMcpServer) RequestSampling(ctx context.Context, request types.CreateMessageRequestParams) (*types.CreateMessageResult, error) {
	return d.server.RequestSampling(ctx, request)
}

// RequestSamplingWithChatMessage implements IMcpServer.
func (d *DestinationBoundMcpServer) RequestSamplingWithChatMessage(ctx context.Context, messages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	return d.server.RequestSamplingWithChatMessage(ctx, messages, options)
}

// Run implements IMcpServer.
func (d *DestinationBoundMcpServer) Run(ctx context.Context) error {
	return d.server.Run(ctx)
}

// SendMessage implements IMcpServer.
func (d *DestinationBoundMcpServer) SendMessage(ctx context.Context, msg transport.IJsonRpcMessage) error {
	msg.SetRelatedTransport(d.transport)
	return d.server.SendMessage(ctx, msg)
}

// SendNotification implements IMcpServer.
func (d *DestinationBoundMcpServer) SendNotification(ctx context.Context, notification transport.JsonRpcNotification) error {
	return d.server.SendNotification(ctx, notification)
}

// SendRequest implements IMcpServer.
func (d *DestinationBoundMcpServer) SendRequest(ctx context.Context, req *transport.JsonRpcRequest) (*transport.JsonRpcResponse, error) {
	req.SetRelatedTransport(d.transport)
	return d.server.SendRequest(ctx, req)
}
