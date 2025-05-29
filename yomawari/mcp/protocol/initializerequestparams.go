package protocol

type InitializeRequestParams struct {
	RequestParams   `json:",inline"`
	ProtocolVersion string              `json:"protocolVersion"`
	Capabilities    *ClientCapabilities `json:"capabilities"`
	ClientInfo      Implementation      `json:"clientInfo"`
}
