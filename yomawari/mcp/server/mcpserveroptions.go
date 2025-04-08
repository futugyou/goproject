package server

import (
	"context"
	"time"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type McpServerOptions struct {
	ServerInfo            types.Implementation
	Capabilities          *ServerCapabilities
	ProtocolVersion       string        // "2024-11-05"
	InitializationTimeout time.Duration //  60 sec.
	ServerInstructions    string
	GetCompletionHandler  func(context.Context, RequestContext[*types.CompleteRequestParams]) types.CompleteResult
}
