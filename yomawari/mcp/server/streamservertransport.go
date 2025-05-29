package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/futugyou/yomawari/mcp/logging"
	"github.com/futugyou/yomawari/mcp/protocol"
)

var _ protocol.ITransport = (*StreamServerTransport)(nil)
var (
	newlineBytes = []byte("\n")
)

// StreamServerTransport provides a server transport implemented around a pair of input/output streams.
type StreamServerTransport struct {
	*protocol.TransportBase
	logger logging.Logger

	inputReader  *bufio.Reader
	outputStream io.Writer
	endpointName string

	sendLock    sync.Mutex
	shutdownCtx context.Context
	cancelFunc  context.CancelFunc

	readLoopCompleted chan struct{}
	disposed          bool
	connectedMutex    sync.RWMutex
}

// GetTransportKind implements ITransport.
func (t *StreamServerTransport) GetTransportKind() protocol.TransportKind {
	return protocol.TransportKindStream
}

// MessageReader implements ITransport.
// Subtle: this method shadows the method (TransportBase).MessageReader of StreamServerTransport.TransportBase.
func (t *StreamServerTransport) MessageReader() <-chan protocol.IJsonRpcMessage {
	panic("unimplemented")
}

// NewStreamServerTransport creates a new StreamServerTransport with explicit input/output streams.
func NewStreamServerTransport(inputStream io.Reader, outputStream io.Writer, serverName string, logger logging.Logger) *StreamServerTransport {
	if inputStream == nil {
		panic("inputStream cannot be nil")
	}
	if outputStream == nil {
		panic("outputStream cannot be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	t := &StreamServerTransport{
		TransportBase:     protocol.ServerTransportBase(),
		logger:            logger,
		inputReader:       bufio.NewReader(inputStream),
		outputStream:      outputStream,
		shutdownCtx:       ctx,
		cancelFunc:        cancel,
		readLoopCompleted: make(chan struct{}),
	}

	if serverName != "" {
		t.endpointName = "Server (stream) (" + serverName + ")"
	} else {
		t.endpointName = "Server (stream)"
	}

	t.SetConnected(true)
	go t.readMessages()

	return t
}

// SendMessage sends a JSON-RPC message through the
func (t *StreamServerTransport) SendMessage(ctx context.Context, message protocol.IJsonRpcMessage) error {
	if !t.IsConnected() {
		t.logger.TransportNotConnected(t.endpointName)
		return fmt.Errorf("transport is not connected")
	}

	t.sendLock.Lock()
	defer t.sendLock.Unlock()

	messageID := "(no id)"
	if msgWithID, ok := message.(protocol.IJsonRpcMessageWithId); ok {
		messageWithId := msgWithID.GetId()
		messageID = messageWithId.String()
	}

	t.logger.TransportSendingMessage(t.endpointName, messageID)

	data, err := json.Marshal(message)
	if err != nil {
		t.logger.TransportSendFailed(t.endpointName, messageID, err)
		return fmt.Errorf("failed to marshal message")
	}

	if _, err := t.outputStream.Write(data); err != nil {
		t.logger.TransportSendFailed(t.endpointName, messageID, err)
		return fmt.Errorf("failed to write message")
	}

	if _, err := t.outputStream.Write(newlineBytes); err != nil {
		t.logger.TransportSendFailed(t.endpointName, messageID, err)
		return fmt.Errorf("failed to write newline")
	}

	t.logger.TransportSentMessage(t.endpointName, messageID)
	return nil
}

func (t *StreamServerTransport) readMessages() {
	defer close(t.readLoopCompleted)

	t.logger.TransportEnteringReadMessagesLoop(t.endpointName)

	for {
		select {
		case <-t.shutdownCtx.Done():
			t.logger.TransportReadMessagesCancelled(t.endpointName)
			return
		default:
			t.logger.TransportWaitingForMessage(t.endpointName)

			line, err := t.inputReader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					t.logger.TransportEndOfStream(t.endpointName)
				} else {
					t.logger.TransportReadMessagesFailed(t.endpointName, err)
				}
				t.SetConnected(false)
				return
			}

			if len(line) == 0 {
				continue
			}

			t.logger.TransportReceivedMessage(t.endpointName, string(line))

			var message protocol.IJsonRpcMessage
			if m, err := protocol.UnmarshalJsonRpcMessage(line); err != nil {
				t.logger.TransportMessageParseFailed(t.endpointName, string(line), err)
				continue
			} else {
				message = m
			}

			messageID := "(no id)"
			if msgWithID, ok := message.(protocol.IJsonRpcMessageWithId); ok {
				messageWithId := msgWithID.GetId()
				messageID = messageWithId.String()
			}

			t.logger.TransportReceivedMessageParsed(t.endpointName, messageID)

			t.WriteMessage(t.shutdownCtx, message)
			t.logger.TransportMessageWritten(t.endpointName, messageID)
		}
	}
}

// Close cleans up the transport resources.
func (t *StreamServerTransport) Close() error {
	if t.disposed {
		return nil
	}
	t.disposed = true

	t.logger.TransportCleaningUp(t.endpointName)

	// Signal to the stdin reading loop to stop
	t.cancelFunc()

	// In Go, we don't have direct control over closing the underlying file descriptors
	// like in C#. The caller should close the streams if needed.

	// Wait for the read loop to complete
	t.logger.TransportWaitingForReadTask(t.endpointName)

	select {
	case <-t.readLoopCompleted:
		t.logger.TransportReadTaskCleanedUp(t.endpointName)
	case <-time.After(5 * time.Second):
		t.logger.TransportCleanupReadTaskTimeout(t.endpointName)
	}

	t.SetConnected(false)
	t.logger.TransportCleanedUp(t.endpointName)
	return nil
}

func (t *StreamServerTransport) IsConnected() bool {
	t.connectedMutex.RLock()
	defer t.connectedMutex.RUnlock()
	return t.TransportBase.IsConnected()
}

func (t *StreamServerTransport) SetConnected(connected bool) {
	t.connectedMutex.Lock()
	defer t.connectedMutex.Unlock()
	t.TransportBase.SetConnected(connected)
}
