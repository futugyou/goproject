package server

import (
	"time"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type McpServerOptions struct {
	ServerInfo            types.Implementation
	Capabilities          *ServerCapabilities
	ProtocolVersion       string        // "2024-11-05"
	InitializationTimeout time.Duration //  60 sec.
	ServerInstructions    string
	ScopeRequests         bool
}

func NewMcpServerOptions() *McpServerOptions {
	return &McpServerOptions{
		ServerInfo:            types.Implementation{},
		Capabilities:          &ServerCapabilities{},
		ProtocolVersion:       "2024-11-05",
		InitializationTimeout: time.Duration(60) * time.Second,
		ServerInstructions:    "",
		ScopeRequests:         true,
	}
}
