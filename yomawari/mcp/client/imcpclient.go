package client

import (
	"context"
	"net/url"

	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/mcp/protocol"
	"github.com/futugyou/yomawari/mcp/server"
	"github.com/futugyou/yomawari/mcp/shared"
)

type IMcpClient interface {
	shared.IMcpEndpoint
	GetServerCapabilities() *server.ServerCapabilities
	GetServerInfo() *protocol.Implementation
	GetServerInstructions() *string
	Ping(ctx context.Context) error
	ListTools(ctx context.Context) ([]McpClientTool, error)
	CallTool(ctx context.Context, toolName string, arguments map[string]interface{}, reporter protocol.IProgressReporter) (*protocol.CallToolResult, error)
	GetPrompt(ctx context.Context, name string, arguments map[string]interface{}) (*protocol.GetPromptResult, error)
	EnumerateTools(ctx context.Context) (<-chan McpClientTool, <-chan error)
	ListPrompts(ctx context.Context, client IMcpClient) ([]McpClientPrompt, error)
	EnumeratePrompts(ctx context.Context, client IMcpClient) (<-chan McpClientPrompt, <-chan error)
	ListResourceTemplates(ctx context.Context, client IMcpClient) ([]McpClientResourceTemplate, error)
	EnumerateResourceTemplates(ctx context.Context, client IMcpClient) (<-chan McpClientResourceTemplate, <-chan error)
	ListResources(ctx context.Context, client IMcpClient) ([]McpClientResource, error)
	EnumerateResources(ctx context.Context, client IMcpClient) (<-chan McpClientResource, <-chan error)
	ReadResource(ctx context.Context, uri string) (*protocol.ReadResourceResult, error)
	ReadResourceWithUri(ctx context.Context, uri url.URL) (*protocol.ReadResourceResult, error)
	ReadResourceWithUriAndArguments(ctx context.Context, uriTemplate string, arguments map[string]interface{}) (*protocol.ReadResourceResult, error)
	Complete(ctx context.Context, reference protocol.Reference, argumentName string, argumentValue string) (*protocol.CompleteResult, error)
	SubscribeToResource(ctx context.Context, uri string) error
	SubscribeToResourceWithUri(ctx context.Context, uri url.URL) error
	UnsubscribeFromResource(ctx context.Context, uri string) error
	UnsubscribeFromResourceWithUri(ctx context.Context, uri url.URL) error
	SetLoggingLevel(ctx context.Context, level protocol.LoggingLevel) error
	SetLoggingLevelWithLogLevel(ctx context.Context, level logger.LogLevel) error
}
