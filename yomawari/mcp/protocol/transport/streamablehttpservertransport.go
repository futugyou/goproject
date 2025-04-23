package transport

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

var _ ITransport = (*StreamableHttpServerTransport)(nil)

type StreamableHttpServerTransport struct {
	sseWriter         *SseWriter
	incomingChannel   chan messages.IJsonRpcMessage
	ctx               context.Context
	cancelFunc        context.CancelFunc
	getRequestStarted int32
}

func NewStreamableHttpServerTransport() *StreamableHttpServerTransport {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &StreamableHttpServerTransport{
		sseWriter:         NewSseWriter(""),
		incomingChannel:   make(chan messages.IJsonRpcMessage),
		ctx:               ctx,
		cancelFunc:        cancelFunc,
		getRequestStarted: 0,
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

func (s *StreamableHttpServerTransport) HandleGetRequest(ctx context.Context, sseResponseStream io.Writer) error {
	if atomic.SwapInt32(&s.getRequestStarted, 1) != 0 {
		return fmt.Errorf("session resumption is not yet supported. Please start a new session")
	}
	ctx, _ = mergeContexts(s.ctx, ctx)
	resultCh := s.sseWriter.WriteAll(ctx, sseResponseStream)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resultCh:
		return err
	}
}

func (s *StreamableHttpServerTransport) HandlePostRequest(ctx context.Context, httpBodies *DuplexPipe) (bool, error) {
	ctx, _ = mergeContexts(s.ctx, ctx)
	postTransport := NewStreamableHttpPostTransport(s.incomingChannel, httpBodies)
	return postTransport.Run(ctx)
}

func mergeContexts(ctx1, ctx2 context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ctx1.Done():
			cancel()
		case <-ctx2.Done():
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}
