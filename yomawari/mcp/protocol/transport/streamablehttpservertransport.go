package transport

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"
)

var _ ITransport = (*StreamableHttpServerTransport)(nil)

type StreamableHttpServerTransport struct {
	sseWriter         *SseWriter
	Stateless         bool
	InitializeRequest InitializeRequestParams
	incomingChannel   chan IJsonRpcMessage
	ctx               context.Context
	cancelFunc        context.CancelFunc
	getRequestStarted int32
}

func NewStreamableHttpServerTransport() *StreamableHttpServerTransport {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &StreamableHttpServerTransport{
		sseWriter:         NewSseWriter(""),
		incomingChannel:   make(chan IJsonRpcMessage),
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
func (s *StreamableHttpServerTransport) MessageReader() <-chan IJsonRpcMessage {
	return s.incomingChannel
}

// SendMessage implements ITransport.
func (s *StreamableHttpServerTransport) SendMessage(ctx context.Context, message IJsonRpcMessage) error {
	if s.Stateless {
		return fmt.Errorf("stateless mode is not supported for GET requests")
	}

	return s.sseWriter.SendMessage(ctx, message)
}

func (s *StreamableHttpServerTransport) HandleGetRequest(ctx context.Context, sseResponseStream io.Writer) error {
	if s.Stateless {
		return fmt.Errorf("stateless mode is not supported for GET requests")
	}

	if atomic.SwapInt32(&s.getRequestStarted, 1) != 0 {
		return fmt.Errorf("session resumption is not yet supported. Please start a new session")
	}

	// We do not need to reference ctx like in HandlePostRequest, because the session ending completes the sseWriter gracefully.
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
	postTransport := NewStreamableHttpPostTransport(s, httpBodies)
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
