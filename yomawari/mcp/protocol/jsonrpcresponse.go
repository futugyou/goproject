package protocol

import "encoding/json"

type JsonRpcResponse struct {
	JsonRpc          string          `json:"jsonrpc"`
	Result           json.RawMessage `json:"result,omitempty"`
	Id               *RequestId      `json:"id"`
	RelatedTransport ITransport      `json:"-"`
}

// GetRelatedTransport implements IJsonRpcMessageWithId.
func (j *JsonRpcResponse) GetRelatedTransport() ITransport {
	return j.RelatedTransport
}

// SetRelatedTransport implements IJsonRpcMessageWithId.
func (j *JsonRpcResponse) SetRelatedTransport(transport ITransport) {
	j.RelatedTransport = transport
}

func NewJsonRpcResponse(id *RequestId, result json.RawMessage) *JsonRpcResponse {
	return &JsonRpcResponse{
		JsonRpc: "2.0",
		Result:  result,
		Id:      id,
	}
}

func NewJsonRpcResponseWithTransport(id *RequestId, result json.RawMessage, transport ITransport) *JsonRpcResponse {
	return &JsonRpcResponse{
		JsonRpc:          "2.0",
		Result:           result,
		Id:               id,
		RelatedTransport: transport,
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
