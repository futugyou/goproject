package server

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/mcp/protocol/types"
	"github.com/futugyou/yomawari/mcp/shared"
)

type IMcpServer interface {
	shared.IMcpEndpoint
	GetClientCapabilities() *transport.ClientCapabilities
	GetClientInfo() *types.Implementation
	GetMcpServerOptions() *McpServerOptions
	Run(ctx context.Context) error
	RequestSampling(ctx context.Context, request types.CreateMessageRequestParams) (*types.CreateMessageResult, error)
	RequestSamplingWithChatMessage(ctx context.Context, messages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error)
	AsSamplingChatClient() (chatcompletion.IChatClient, error)
	RequestRoots(ctx context.Context, request types.ListRootsRequestParams) (*types.ListRootsResult, error)
}
