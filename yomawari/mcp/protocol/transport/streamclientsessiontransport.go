package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/futugyou/yomawari/mcp/logging"
	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

var _ ITransport = (*StreamClientSessionTransport)(nil)

// StreamClientSessionTransport represents an active client session transport.
type StreamClientSessionTransport struct {
	TransportBase
	logger logging.Logger

	serverOutput *bufio.Reader
	serverInput  io.Writer
	EndpointName string

	sendLock          sync.Mutex
	disposed          bool
	connectedMu       sync.RWMutex
	shutdownCtx       context.Context
	cancelFunc        context.CancelFunc
	readLoopCompleted chan struct{}
}

// NewStreamClientSessionTransport creates a new StreamClientSessionTransport.
func NewStreamClientSessionTransport(
	serverInput io.Writer,
	serverOutput io.Reader,
	endpointName string,
	logger logging.Logger,
) *StreamClientSessionTransport {
	ctx, cancel := context.WithCancel(context.Background())
	t := &StreamClientSessionTransport{
		TransportBase: TransportBase{
			messageChannel: make(chan messages.IJsonRpcMessage),
			isConnected:    false,
		},
		logger:            logger,
		serverOutput:      bufio.NewReader(serverOutput),
		serverInput:       serverInput,
		EndpointName:      endpointName,
		shutdownCtx:       ctx,
		cancelFunc:        cancel,
		readLoopCompleted: make(chan struct{}),
	}
	t.SetConnected(true)
	go t.readMessages()
	return t
}

func (t *StreamClientSessionTransport) IsConnected() bool {
	t.connectedMu.RLock()
	defer t.connectedMu.RUnlock()
	return t.isConnected
}

func (t *StreamClientSessionTransport) SetConnected(connected bool) {
	t.connectedMu.Lock()
	defer t.connectedMu.Unlock()
	t.isConnected = connected
}

// Close closes the transport and releases resources.
func (t *StreamClientSessionTransport) Close() error {
	if t.disposed {
		return nil
	}
	t.disposed = true

	t.logger.TransportCleaningUp(t.EndpointName)

	// Signal to the stdin reading loop to stop
	t.cancelFunc()

	// Wait for the read loop to complete
	t.logger.TransportWaitingForReadTask(t.EndpointName)

	select {
	case <-t.readLoopCompleted:
		t.logger.TransportReadTaskCleanedUp(t.EndpointName)
	case <-time.After(5 * time.Second):
		t.logger.TransportCleanupReadTaskTimeout(t.EndpointName)
	}

	t.SetConnected(false)
	t.logger.TransportCleanedUp(t.EndpointName)
	return nil
}

func (t *StreamClientSessionTransport) readMessages() {
	defer close(t.readLoopCompleted)

	t.logger.TransportEnteringReadMessagesLoop(t.EndpointName)

	for {
		select {
		case <-t.shutdownCtx.Done():
			t.logger.TransportReadMessagesCancelled(t.EndpointName)
			return
		default:
			t.logger.TransportWaitingForMessage(t.EndpointName)

			line, err := t.serverOutput.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					t.logger.TransportEndOfStream(t.EndpointName)
				} else {
					t.logger.TransportReadMessagesFailed(t.EndpointName, err)
				}
				t.SetConnected(false)
				return
			}

			if len(line) == 0 {
				continue
			}

			t.logger.TransportReceivedMessage(t.EndpointName, string(line))

			var message messages.IJsonRpcMessage
			if m, err := messages.UnmarshalJsonRpcMessage(line); err != nil {
				t.logger.TransportMessageParseFailed(t.EndpointName, string(line), err)
				continue
			} else {
				message = m
			}

			messageID := "(no id)"
			if msgWithID, ok := message.(messages.IJsonRpcMessageWithId); ok {
				messageWithId := msgWithID.GetId()
				messageID = messageWithId.String()
			}

			t.logger.TransportReceivedMessageParsed(t.EndpointName, messageID)

			t.WriteMessage(t.shutdownCtx, message)
			t.logger.TransportMessageWritten(t.EndpointName, messageID)
		}
	}
}

// SendMessageAsync implements ITransport.
func (t *StreamClientSessionTransport) SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error {
	if !t.IsConnected() {
		t.logger.TransportNotConnected(t.EndpointName)
		return fmt.Errorf("transport is not connected")
	}

	t.sendLock.Lock()
	defer t.sendLock.Unlock()

	messageID := "(no id)"
	if msgWithID, ok := message.(messages.IJsonRpcMessageWithId); ok {
		messageWithId := msgWithID.GetId()
		messageID = messageWithId.String()
	}

	t.logger.TransportSendingMessage(t.EndpointName, messageID)

	data, err := json.Marshal(message)
	if err != nil {
		t.logger.TransportSendFailed(t.EndpointName, messageID, err)
		return fmt.Errorf("failed to marshal message")
	}

	if _, err := t.serverInput.Write(data); err != nil {
		t.logger.TransportSendFailed(t.EndpointName, messageID, err)
		return fmt.Errorf("failed to write message")
	}

	t.logger.TransportSentMessage(t.EndpointName, messageID)
	return nil
}
