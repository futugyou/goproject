package protocol

import (
	"context"
	"fmt"
	"io"
	"sync"
)

var _ ITransport = (*SseResponseStreamTransport)(nil)

type SseResponseStreamTransport struct {
	sseResponseStream io.Writer
	messageEndpoint   string

	incomingChannel chan IJsonRpcMessage

	isConnected bool
	mu          sync.Mutex
	sseWriter   *SseWriter
}

func NewSseResponseStreamTransport(sseResponseStream io.Writer, messageEndpoint string) *SseResponseStreamTransport {
	if messageEndpoint == "" {
		messageEndpoint = "/message"
	}

	return &SseResponseStreamTransport{
		sseResponseStream: sseResponseStream,
		messageEndpoint:   messageEndpoint,
		incomingChannel:   make(chan IJsonRpcMessage, 1), // Buffered channel
		sseWriter:         NewSseWriter(messageEndpoint),
	}
}

// Close implements I
func (t *SseResponseStreamTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.isConnected {
		return nil // Already closed
	}

	t.isConnected = false

	// Close channels
	close(t.incomingChannel)

	t.sseWriter.Dispose()
	return nil
}

// MessageReader implements I
func (s *SseResponseStreamTransport) MessageReader() <-chan IJsonRpcMessage {
	return s.incomingChannel
}

// SendMessage implements I
func (t *SseResponseStreamTransport) SendMessage(ctx context.Context, message IJsonRpcMessage) error {
	return t.sseWriter.SendMessage(ctx, message)
}

// Run starts the transport and writes the JSON-RPC messages to the SSE response stream
func (t *SseResponseStreamTransport) Run(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.isConnected = true

	errCh := t.sseWriter.WriteAll(ctx, t.sseResponseStream)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

// OnMessageReceived handles incoming JSON-RPC messages
func (t *SseResponseStreamTransport) OnMessageReceived(ctx context.Context, message IJsonRpcMessage) error {
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
