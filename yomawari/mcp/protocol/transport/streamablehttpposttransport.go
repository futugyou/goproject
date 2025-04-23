package transport

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

var _ ITransport = (*StreamableHttpPostTransport)(nil)

type StreamableHttpPostTransport struct {
	httpBodies      *DuplexPipe
	incomingChannel chan messages.IJsonRpcMessage
	sseWriter       *SseWriter
	pendingRequests map[messages.RequestId]struct{}
}

func NewStreamableHttpPostTransport(incomingChannel chan messages.IJsonRpcMessage, httpBodies *DuplexPipe) *StreamableHttpPostTransport {

	return &StreamableHttpPostTransport{
		httpBodies:      httpBodies,
		incomingChannel: incomingChannel,
		sseWriter:       NewSseWriter(""),
		pendingRequests: make(map[messages.RequestId]struct{}),
	}
}

// Close implements ITransport.
func (s *StreamableHttpPostTransport) Close() error {
	s.sseWriter.Dispose()
	return nil
}

// MessageReader implements ITransport.
func (s *StreamableHttpPostTransport) MessageReader() <-chan messages.IJsonRpcMessage {
	panic("JsonRpcMessage.RelatedTransport should only be used for sending messages.")
}

// SendMessage implements ITransport.
func (s *StreamableHttpPostTransport) SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error {
	return s.sseWriter.SendMessage(ctx, message)
}
