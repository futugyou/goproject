package transport

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

var _ ITransport = (*StreamableHttpServerTransport)(nil)

type StreamableHttpServerTransport struct {
	sseWriter         *SseWriter
	incomingChannel   chan messages.IJsonRpcMessage
	ctx               context.Context
	cancelFunc        context.CancelFunc
	getRequestStarted int
}

func NewStreamableHttpServerTransport() *StreamableHttpServerTransport {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &StreamableHttpServerTransport{
		sseWriter:       NewSseWriter(""),
		incomingChannel: make(chan messages.IJsonRpcMessage),
		ctx:             ctx,
		cancelFunc:      cancelFunc,
	}
}

// Close implements ITransport.
func (s *StreamableHttpServerTransport) Close() error {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	s.sseWriter.Dispose()
	return nil
}

// MessageReader implements ITransport.
func (s *StreamableHttpServerTransport) MessageReader() <-chan messages.IJsonRpcMessage {
	return s.incomingChannel
}

// SendMessage implements ITransport.
func (s *StreamableHttpServerTransport) SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error {
	return s.sseWriter.SendMessage(ctx, message)
}
