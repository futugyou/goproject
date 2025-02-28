package contents

// FunctionResultContent represents the result of a function call.
type FunctionResultContent struct {
	AIContent `json:",inline"`
	CallId    string      `json:"callId,omitempty"`
	Result    interface{} `json:"result,omitempty"`
	Error     error       `json:"-"`
}
