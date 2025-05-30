package server

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/mcp/protocol"
	"github.com/futugyou/yomawari/mcp/shared"
)

type IMcpServer interface {
	shared.IMcpEndpoint
	GetClientCapabilities() *protocol.ClientCapabilities
	GetClientInfo() *protocol.Implementation
	GetMcpServerOptions() *McpServerOptions
	Run(ctx context.Context) error
	Sample(ctx context.Context, request protocol.CreateMessageRequestParams) (*protocol.CreateMessageResult, error)
	Elicit(ctx context.Context, request protocol.ElicitRequestParams) (*protocol.ElicitResult, error)
	SampleWithChatMessage(ctx context.Context, messages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error)
	AsSamplingChatClient() (chatcompletion.IChatClient, error)
	RequestRoots(ctx context.Context, request protocol.ListRootsRequestParams) (*protocol.ListRootsResult, error)
}
