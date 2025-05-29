package server

import "github.com/futugyou/yomawari/mcp/protocol"

type InitializeResult struct {
	ProtocolVersion string                  `json:"protocolVersion"`
	Capabilities    ServerCapabilities      `json:"capabilities"`
	ServerInfo      protocol.Implementation `json:"serverInfo"`
	Instructions    string                  `json:"instructions"`
}
