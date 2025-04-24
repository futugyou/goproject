package transport

type JsonRpcError struct {
	JsonRpc string              `json:"jsonrpc"`
	Id      *RequestId          `json:"id"`
	Error   *JsonRpcErrorDetail `json:"error"`
}

func NewJsonRpcError(id *RequestId, code int, message string, data any) *JsonRpcError {
	return &JsonRpcError{
		JsonRpc: "2.0",
		Id:      id,
		Error: &JsonRpcErrorDetail{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// GetId implements IJsonRpcMessageWithId.
func (j *JsonRpcError) GetId() *RequestId {
	return j.Id
}

// GetJsonRpc implements IJsonRpcMessageWithId.
func (j *JsonRpcError) GetJsonRpc() string {
	return "2.0"
}

type JsonRpcErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

var _ IJsonRpcMessageWithId = (*JsonRpcError)(nil)
