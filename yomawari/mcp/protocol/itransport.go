package protocol

import (
	"context"
)

type ITransport interface {
	MessageReader() <-chan IJsonRpcMessage
	SendMessage(ctx context.Context, message IJsonRpcMessage) error
	Close() error
	GetTransportKind() TransportKind
}

const TransportTypesStdIo string = "stdio"
const TransportTypesSse string = "sse"

type TransportKind string

var TransportKindUnknown TransportKind = "unknownTransport"
var TransportKindStdio TransportKind = "stdio"
var TransportKindStream TransportKind = "stream"
var TransportKindSse TransportKind = "sse"
var TransportKindHttp TransportKind = "http"
