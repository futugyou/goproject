package client

import (
	"time"

	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type McpClientOptions struct {
	ClientInfo            *types.Implementation
	Capabilities          *transport.ClientCapabilities
	ProtocolVersion       string
	InitializationTimeout time.Duration
}

func NewMcpClientOptions() *McpClientOptions {
	return &McpClientOptions{
		ClientInfo:            &types.Implementation{},
		Capabilities:          &transport.ClientCapabilities{},
		ProtocolVersion:       "2024-11-05",
		InitializationTimeout: time.Duration(60) * time.Second,
	}
}
