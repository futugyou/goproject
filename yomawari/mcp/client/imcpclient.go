package client

import (
	"context"
	"net/url"

	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/mcp/protocol/types"
	"github.com/futugyou/yomawari/mcp/server"
	"github.com/futugyou/yomawari/mcp/shared"
)

type IMcpClient interface {
	shared.IMcpEndpoint
	GetServerCapabilities() *server.ServerCapabilities
	GetServerInfo() *types.Implementation
	GetServerInstructions() *string
	Ping(ctx context.Context) error
	ListTools(ctx context.Context) ([]McpClientTool, error)
	CallTool(ctx context.Context, toolName string, arguments map[string]interface{}, reporter transport.IProgressReporter) (*types.CallToolResult, error)
	GetPrompt(ctx context.Context, name string, arguments map[string]interface{}) (*types.GetPromptResult, error)
	EnumerateTools(ctx context.Context) (<-chan McpClientTool, <-chan error)
	ListPrompts(ctx context.Context, client IMcpClient) ([]McpClientPrompt, error)
	EnumeratePrompts(ctx context.Context, client IMcpClient) (<-chan McpClientPrompt, <-chan error)
	ListResourceTemplates(ctx context.Context, client IMcpClient) ([]McpClientResourceTemplate, error)
	EnumerateResourceTemplates(ctx context.Context, client IMcpClient) (<-chan McpClientResourceTemplate, <-chan error)
	ListResources(ctx context.Context, client IMcpClient) ([]McpClientResource, error)
	EnumerateResources(ctx context.Context, client IMcpClient) (<-chan McpClientResource, <-chan error)
	ReadResource(ctx context.Context, uri string) (*types.ReadResourceResult, error)
	ReadResourceWithUri(ctx context.Context, uri url.URL) (*types.ReadResourceResult, error)
	ReadResourceWithUriAndArguments(ctx context.Context, uriTemplate string, arguments map[string]interface{}) (*types.ReadResourceResult, error)
	Complete(ctx context.Context, reference types.Reference, argumentName string, argumentValue string) (*types.CompleteResult, error)
	SubscribeToResource(ctx context.Context, uri string) error
	SubscribeToResourceWithUri(ctx context.Context, uri url.URL) error
	UnsubscribeFromResource(ctx context.Context, uri string) error
	UnsubscribeFromResourceWithUri(ctx context.Context, uri url.URL) error
	SetLoggingLevel(ctx context.Context, level types.LoggingLevel) error
	SetLoggingLevelWithLogLevel(ctx context.Context, level logger.LogLevel) error
}
