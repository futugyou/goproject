package transport

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

var _ ITransport = (*SseResponseStreamTransport)(nil)

type SseResponseStreamTransport struct {
	sseResponseStream io.Writer
	messageEndpoint   string

	incomingChannel    chan messages.IJsonRpcMessage
	outgoingSseChannel chan SseItem

	sseWriteTaskCancel context.CancelFunc
	sseWriteTaskWg     sync.WaitGroup

	isConnected bool
	mu          sync.Mutex
}

func NewSseResponseStreamTransport(sseResponseStream io.Writer, messageEndpoint string) *SseResponseStreamTransport {
	if messageEndpoint == "" {
		messageEndpoint = "/message"
	}

	return &SseResponseStreamTransport{
		sseResponseStream:  sseResponseStream,
		messageEndpoint:    messageEndpoint,
		incomingChannel:    make(chan messages.IJsonRpcMessage, 1), // Buffered channel
		outgoingSseChannel: make(chan SseItem, 1),                  // Buffered channel
	}
}

// Close implements ITransport.
func (t *SseResponseStreamTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.isConnected {
		return nil // Already closed
	}

	t.isConnected = false

	// Cancel the SSE write task
	if t.sseWriteTaskCancel != nil {
		t.sseWriteTaskCancel()
	}

	// Close channels
	close(t.incomingChannel)
	close(t.outgoingSseChannel)

	// Wait for the SSE write task to finish
	t.sseWriteTaskWg.Wait()

	return nil
}

// IsConnected implements ITransport.
func (s *SseResponseStreamTransport) IsConnected() bool {
	return s.isConnected
}

// MessageReader implements ITransport.
func (s *SseResponseStreamTransport) MessageReader() <-chan messages.IJsonRpcMessage {
	return s.incomingChannel
}

// SendMessage implements ITransport.
func (t *SseResponseStreamTransport) SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error {
	t.mu.Lock()
	if !t.isConnected {
		t.mu.Unlock()
		return fmt.Errorf("transport is not connected. Make sure to call Run first")
	}
	t.mu.Unlock()

	data, err := messages.MarshalJsonRpcMessage(message)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case t.outgoingSseChannel <- SseItem{Data: string(data), EventType: "message"}:
		return nil
	}
}

// Run starts the transport and writes the JSON-RPC messages to the SSE response stream
func (t *SseResponseStreamTransport) Run(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.isConnected {
		return nil // Already running
	}

	// Write the initial endpoint event
	select {
	case t.outgoingSseChannel <- SseItem{Data: "", EventType: "endpoint"}:
	default:
		return fmt.Errorf("failed to write endpoint event - channel full")
	}

	t.isConnected = true

	// Create a cancellable context for the SSE write task
	var sseCtx context.Context
	sseCtx, t.sseWriteTaskCancel = context.WithCancel(ctx)

	t.sseWriteTaskWg.Add(1)
	go func() {
		defer t.sseWriteTaskWg.Done()
		t.sseWriteLoop(sseCtx)
	}()

	return nil
}

// sseWriteLoop handles writing SSE items to the response stream
func (t *SseResponseStreamTransport) sseWriteLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-t.outgoingSseChannel:
			if !ok {
				return // Channel closed
			}

			if item.EventType == "endpoint" {
				_, err := t.sseResponseStream.Write([]byte(t.messageEndpoint))
				if err != nil {
					return // Error writing
				}
				continue
			}

			// Write the SSE formatted data
			_, err := t.sseResponseStream.Write([]byte(item.Data))
			if err != nil {
				return // Error writing
			}
		}
	}
}

// OnMessageReceived handles incoming JSON-RPC messages
func (t *SseResponseStreamTransport) OnMessageReceived(ctx context.Context, message messages.IJsonRpcMessage) error {
	t.mu.Lock()
	if !t.isConnected {
		t.mu.Unlock()
		return fmt.Errorf("transport is not connected. Make sure to call Run first")
	}
	t.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case t.incomingChannel <- message:
		return nil
	}
}
