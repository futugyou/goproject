package messages

type JsonRpcNotification struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

func NewJsonRpcNotification(method string, params any) *JsonRpcNotification {
	return &JsonRpcNotification{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

// GetJsonRpc implements IJsonRpcMessage.
func (j *JsonRpcNotification) GetJsonRpc() string {
	return "2.0"
}

var _ IJsonRpcMessage = (*JsonRpcNotification)(nil)
