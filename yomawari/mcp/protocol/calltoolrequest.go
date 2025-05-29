package protocol

// Used by the client to invoke a tool provided by the server.
type CallToolRequest struct {
	// Method corresponds to the JSON schema field "method".
	Method string `json:"method" yaml:"method" mapstructure:"method"` // tools/call

	// Params corresponds to the JSON schema field "params".
	Params CallToolRequestParams `json:"params" yaml:"params" mapstructure:"params"`
}
