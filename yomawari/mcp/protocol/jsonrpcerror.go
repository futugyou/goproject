package protocol

type JsonRpcError struct {
	JsonRpc          string              `json:"jsonrpc"`
	Id               *RequestId          `json:"id"`
	Error            *JsonRpcErrorDetail `json:"error"`
	RelatedTransport ITransport          `json:"-"`
}

// GetRelatedTransport implements IJsonRpcMessageWithId.
func (j *JsonRpcError) GetRelatedTransport() ITransport {
	return j.RelatedTransport
}

// SetRelatedTransport implements IJsonRpcMessageWithId.
func (j *JsonRpcError) SetRelatedTransport(transport ITransport) {
	j.RelatedTransport = transport
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

func NewJsonRpcErrorWithTransport(id *RequestId, code int, message string, data any, transport ITransport) *JsonRpcError {
	return &JsonRpcError{
		JsonRpc:          "2.0",
		Id:               id,
		RelatedTransport: transport,
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
