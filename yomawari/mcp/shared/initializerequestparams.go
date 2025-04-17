package shared

import "github.com/futugyou/yomawari/mcp/protocol/types"

type InitializeRequestParams struct {
	types.RequestParams `json:",inline"`
	ProtocolVersion     string               `json:"protocolVersion"`
	Capabilities        *ClientCapabilities  `json:"capabilities"`
	ClientInfo          types.Implementation `json:"clientInfo"`
}
