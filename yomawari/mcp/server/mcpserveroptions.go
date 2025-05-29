package server

import (
	"time"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type McpServerOptions struct {
	ServerInfo            protocol.Implementation
	Capabilities          *ServerCapabilities
	ProtocolVersion       string        // "2024-11-05"
	InitializationTimeout time.Duration //  60 sec.
	ServerInstructions    string
	ScopeRequests         bool
	KnownClientInfo       *protocol.Implementation
}

func NewMcpServerOptions() *McpServerOptions {
	return &McpServerOptions{
		ServerInfo:            protocol.Implementation{},
		Capabilities:          &ServerCapabilities{},
		ProtocolVersion:       "2024-11-05",
		InitializationTimeout: time.Duration(60) * time.Second,
		ServerInstructions:    "",
		ScopeRequests:         true,
	}
}
