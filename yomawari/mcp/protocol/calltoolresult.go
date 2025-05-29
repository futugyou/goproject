package protocol

// The server's response to a tool call.
//
// Any errors that originate from the tool SHOULD be reported inside the result
// object, with `isError` set to true, _not_ as an MCP protocol-level error
// response. Otherwise, the LLM would not be able to see that an error occurred
// and self-correct.
//
// However, any errors in _finding_ the tool, an error indicating that the
// server does not support tool calls, or any other exceptional conditions,
// should be reported as an MCP error response.
type CallToolResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta map[string]interface{} `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Content corresponds to the JSON schema field "content".
	Content []Content `json:"content" yaml:"content" mapstructure:"content"`

	// Whether the tool call ended in an error.
	//
	// If not set, this is assumed to be false (the call was successful).
	IsError bool `json:"isError,omitempty" yaml:"isError,omitempty" mapstructure:"isError,omitempty"`
}

func NewCallToolResult() *CallToolResult {
	return &CallToolResult{
		Meta:    make(map[string]interface{}),
		Content: make([]Content, 0),
	}
}

func NewCallToolResultWithContents(contents []Content) *CallToolResult {
	return &CallToolResult{
		Meta:    make(map[string]interface{}),
		Content: contents,
	}
}

func NewCallToolResultWithContent(content Content) *CallToolResult {
	return &CallToolResult{
		Meta:    make(map[string]interface{}),
		Content: []Content{content},
	}
}
