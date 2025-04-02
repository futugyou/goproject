package types

type InitializeRequestParams struct {
	RequestParams
	ProtocolVersion string              `json:"protocolVersion"`
	Capabilities    *ClientCapabilities `json:"capabilities"`
	ClientInfo      Implementation      `json:"clientInfo"`
}
