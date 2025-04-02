package server

import "github.com/futugyou/yomawari/mcp/protocol/types"

type InitializeResult struct {
	ProtocolVersion string               `json:"protocolVersion"`
	Capabilities    ServerCapabilities   `json:"capabilities"`
	ServerInfo      types.Implementation `json:"serverInfo"`
	Instructions    string               `json:"instructions"`
}
