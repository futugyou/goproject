package transport

import (
	"context"
)

type ITransport interface {
	MessageReader() <-chan IJsonRpcMessage
	SendMessage(ctx context.Context, message IJsonRpcMessage) error
	Close() error
}

const TransportTypesStdIo string = "stdio"
const TransportTypesSse string = "sse"
