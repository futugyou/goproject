package protocol

import (
	"context"
	"fmt"
)

type TransportBase struct {
	messageChannel chan IJsonRpcMessage
	isConnected    bool
	Name           string
}

// GetTransportKind implements ITransport.
func (t *TransportBase) GetTransportKind() TransportKind {
	panic("unimplemented")
}

func NewTransportBase(name string, messageChannel chan IJsonRpcMessage) *TransportBase {
	if messageChannel == nil {
		messageChannel = make(chan IJsonRpcMessage)
	}
	return &TransportBase{
		Name:           name,
		messageChannel: messageChannel,
		isConnected:    false,
	}
}

func ServerTransportBase() *TransportBase {
	return &TransportBase{
		messageChannel: make(chan IJsonRpcMessage),
		isConnected:    true,
	}
}

func ClientTransportBase() *TransportBase {
	return &TransportBase{
		messageChannel: make(chan IJsonRpcMessage),
		isConnected:    false,
	}
}

// Close implements ITransport.
func (t *TransportBase) Close() error {
	panic("unimplemented")
}

// SendMessageAsync implements ITransport.
func (t *TransportBase) SendMessage(ctx context.Context, message IJsonRpcMessage) error {
	panic("unimplemented")
}

// IsConnected implements ITransport.
func (t *TransportBase) IsConnected() bool {
	return t.isConnected
}

// MessageReader implements ITransport.
func (t *TransportBase) MessageReader() <-chan IJsonRpcMessage {
	return t.messageChannel
}

func (t *TransportBase) WriteMessage(ctx context.Context, message IJsonRpcMessage) error {
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

func MergeContexts(ctx1, ctx2 context.Context) (context.Context, context.CancelFunc) {
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
