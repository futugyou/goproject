package protocol

// The client's response to a sampling/create_message request from the server. The
// client should inform the user before returning the sampled message, to allow
// them to inspect the response (human in the loop) and decide whether to allow the
// server to see it.
type CreateMessageResult struct {
	// This result property is reserved by the protocol to allow clients and servers
	// to attach additional metadata to their responses.
	Meta map[string]interface{} `json:"_meta,omitempty" yaml:"_meta,omitempty" mapstructure:"_meta,omitempty"`

	// Content corresponds to the JSON schema field "content".
	Content Content `json:"content" yaml:"content" mapstructure:"content"`

	// The name of the model that generated the message.
	Model string `json:"model" yaml:"model" mapstructure:"model"`

	// Role corresponds to the JSON schema field "role".
	Role Role `json:"role" yaml:"role" mapstructure:"role"`

	// The reason why sampling stopped, if known.
	StopReason *string `json:"stopReason,omitempty" yaml:"stopReason,omitempty" mapstructure:"stopReason,omitempty"`
}
