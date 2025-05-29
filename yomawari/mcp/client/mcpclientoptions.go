package client

import (
	"time"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type McpClientOptions struct {
	ClientInfo            *protocol.Implementation
	Capabilities          *protocol.ClientCapabilities
	ProtocolVersion       string
	InitializationTimeout time.Duration
}

func NewMcpClientOptions() *McpClientOptions {
	return &McpClientOptions{
		ClientInfo:            &protocol.Implementation{},
		Capabilities:          &protocol.ClientCapabilities{},
		ProtocolVersion:       "2024-11-05",
		InitializationTimeout: time.Duration(60) * time.Second,
	}
}
