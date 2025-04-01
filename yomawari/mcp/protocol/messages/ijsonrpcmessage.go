package messages

type IJsonRpcMessage interface {
	GetJsonRpc() string
}

type IJsonRpcMessageWithId interface {
	IJsonRpcMessage
	GetId() *RequestId
}
