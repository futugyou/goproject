package transport

import "encoding/json"

type JsonRpcResponse struct {
	JsonRpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Id      *RequestId      `json:"id"`
}

func NewJsonRpcResponse(id *RequestId, result json.RawMessage) *JsonRpcResponse {
	return &JsonRpcResponse{
		JsonRpc: "2.0",
		Result:  result,
		Id:      id,
	}
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
