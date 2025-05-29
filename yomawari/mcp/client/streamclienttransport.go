package client

import (
	"context"
	"io"
	"sync"

	"github.com/futugyou/yomawari/mcp/logging"
	"github.com/futugyou/yomawari/mcp/protocol"
)

var _ IClientTransport = (*StreamClientTransport)(nil)

// StreamClientTransport provides a client transport implemented around a pair of input/output streams.
type StreamClientTransport struct {
	serverInput  io.Writer
	serverOutput io.Reader
	logger       logging.Logger
	mu           sync.Mutex
}

// GetName implements IClientTransport.
func (t *StreamClientTransport) GetName() string {
	return "in-memory-stream"
}

// NewStreamClientTransport creates a new StreamClientTransport with explicit input/output streams.
func NewStreamClientTransport(serverInput io.Writer, serverOutput io.Reader, logger logging.Logger) *StreamClientTransport {
	if serverInput == nil {
		panic("serverInput cannot be nil")
	}
	if serverOutput == nil {
		panic("serverOutput cannot be nil")
	}

	return &StreamClientTransport{
		serverInput:  serverInput,
		serverOutput: serverOutput,
		logger:       logger,
	}
}

// Connect creates a new client session transport using the configured streams.
func (t *StreamClientTransport) Connect(ctx context.Context) (protocol.ITransport, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return NewStreamClientSessionTransport(
		t.serverInput,
		t.serverOutput,
		"Client (stream)",
		t.logger,
	), nil
}
