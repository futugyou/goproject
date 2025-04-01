package messages

type JsonRpcResponse struct {
	JsonRpc string     `json:"jsonrpc"`
	Result  any        `json:"result,omitempty"`
	Id      *RequestId `json:"id"`
}

// GetId implements IJsonRpcMessageWithId.
func (j *JsonRpcResponse) GetId() *RequestId {
	return j.Id
}

// GetJsonRpc implements IJsonRpcMessageWithId.
func (j *JsonRpcResponse) GetJsonRpc() string {
	return "2.0"
}

var _ IJsonRpcMessageWithId = (*JsonRpcResponse)(nil)
