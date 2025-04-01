package messages

type JsonRpcNotification struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// GetJsonRpc implements IJsonRpcMessage.
func (j *JsonRpcNotification) GetJsonRpc() string {
	return "2.0"
}

var _ IJsonRpcMessage = (*JsonRpcNotification)(nil)
