package transport

import (
	"context"
	"fmt"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

type TransportBase struct {
	messageChannel chan messages.IJsonRpcMessage
	isConnected    bool
}

func NewTransportBase() *TransportBase {
	return &TransportBase{
		messageChannel: make(chan messages.IJsonRpcMessage),
		isConnected:    true,
	}
}

// Close implements ITransport.
func (t *TransportBase) Close() error {
	panic("unimplemented")
}

// SendMessageAsync implements ITransport.
func (t *TransportBase) SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error {
	panic("unimplemented")
}

// IsConnected implements ITransport.
func (t *TransportBase) IsConnected() bool {
	return t.isConnected
}

// MessageReader implements ITransport.
func (t *TransportBase) MessageReader() <-chan messages.IJsonRpcMessage {
	return t.messageChannel
}

func (t *TransportBase) WriteMessage(ctx context.Context, message messages.IJsonRpcMessage) error {
	if !t.isConnected {
		return fmt.Errorf("transport is not connected")
	}

	select {
	case t.messageChannel <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *TransportBase) SetConnected(isConnected bool) {
	if t.isConnected == isConnected {
		return
	}

	t.isConnected = isConnected
	if !t.isConnected {
		close(t.messageChannel)
	}
}

var _ ITransport = (*TransportBase)(nil)
