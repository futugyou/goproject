package messages

type JsonRpcRequest struct {
	JsonRpc string     `json:"jsonrpc"`
	Method  string     `json:"method"`
	Params  any        `json:"params,omitempty"`
	Id      *RequestId `json:"id"`
}

// GetId implements IJsonRpcMessageWithId.
func (j *JsonRpcRequest) GetId() *RequestId {
	return j.Id
}

// GetJsonRpc implements IJsonRpcMessageWithId.
func (j *JsonRpcRequest) GetJsonRpc() string {
	return "2.0"
}

var _ IJsonRpcMessageWithId = (*JsonRpcRequest)(nil)
