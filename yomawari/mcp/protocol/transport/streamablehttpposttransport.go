package transport

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
	"github.com/futugyou/yomawari/runtime/sse"
)

var _ ITransport = (*StreamableHttpPostTransport)(nil)

type StreamableHttpPostTransport struct {
	httpBodies      *DuplexPipe
	incomingChannel chan messages.IJsonRpcMessage
	sseWriter       *SseWriter
	pendingRequests map[messages.RequestId]struct{}
	mu              sync.Mutex
}

func NewStreamableHttpPostTransport(incomingChannel chan messages.IJsonRpcMessage, httpBodies *DuplexPipe) *StreamableHttpPostTransport {

	return &StreamableHttpPostTransport{
		httpBodies:      httpBodies,
		incomingChannel: incomingChannel,
		sseWriter:       NewSseWriter(""),
		pendingRequests: make(map[messages.RequestId]struct{}),
	}
}

// Close implements ITransport.
func (s *StreamableHttpPostTransport) Close() error {
	s.sseWriter.Dispose()
	return nil
}

// MessageReader implements ITransport.
func (s *StreamableHttpPostTransport) MessageReader() <-chan messages.IJsonRpcMessage {
	panic("JsonRpcMessage.RelatedTransport should only be used for sending messages.")
}

// SendMessage implements ITransport.
func (s *StreamableHttpPostTransport) SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error {
	return s.sseWriter.SendMessage(ctx, message)
}

func (s *StreamableHttpPostTransport) Run(ctx context.Context) (bool, error) {
	if s.incomingChannel != nil {
		if err := s.onPostBodyReceived(ctx); err != nil {
			return false, err
		}
	}

	s.mu.Lock()
	noRequests := len(s.pendingRequests) == 0
	s.mu.Unlock()

	if noRequests {
		return false, nil
	}

	s.sseWriter.MessageFilter = s.stopOnFinalResponseFilter
	s.sseWriter.WriteAll(ctx, s.httpBodies.Output)
	return true, nil
}

func (s *StreamableHttpPostTransport) onPostBodyReceived(ctx context.Context) error {
	var buf bytes.Buffer
	tee := io.TeeReader(s.httpBodies.Input, &buf)

	reader := bufio.NewReader(tee)
	first, err := reader.Peek(1)
	if err != nil {
		return err
	}

	rest := io.MultiReader(reader, s.httpBodies.Input)

	if first[0] == '[' {
		decoder := json.NewDecoder(rest)
		if token, err := decoder.Token(); err != nil || token != json.Delim('[') {
			return fmt.Errorf("invalid JSON array")
		}
		for decoder.More() {
			var raw json.RawMessage
			if err := decoder.Decode(&raw); err != nil {
				return err
			}
			msg, err := messages.UnmarshalJsonRpcMessage([]byte(raw))
			if err != nil {
				return err
			}
			if err := s.onMessageReceived(ctx, msg); err != nil {
				return err
			}
		}
	} else {
		data, err := io.ReadAll(rest)
		if err != nil {
			return err
		}
		msg, err := messages.UnmarshalJsonRpcMessage(data)
		if err != nil {
			return err
		}
		if err := s.onMessageReceived(ctx, msg); err != nil {
			return err
		}
	}

	return nil
}

func (s *StreamableHttpPostTransport) onMessageReceived(ctx context.Context, msg messages.IJsonRpcMessage) error {
	if msg == nil {
		return fmt.Errorf("received invalid null message")
	}

	if condition, ok := msg.(*messages.JsonRpcRequest); ok {
		s.mu.Lock()
		s.pendingRequests[*condition.GetId()] = struct{}{}
		s.mu.Unlock()
	}

	// TODO: move all messages package to transport
	// msg.RelatedTransport = this;

	s.incomingChannel <- msg
	if s.incomingChannel == nil {
		return fmt.Errorf("incoming channel is nil")
	}

	select {
	case s.incomingChannel <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *StreamableHttpPostTransport) stopOnFinalResponseFilter(ctx context.Context, mesg chan sse.SseItem[messages.IJsonRpcMessage]) chan sse.SseItem[messages.IJsonRpcMessage] {
	output := make(chan sse.SseItem[messages.IJsonRpcMessage])
	go func() {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-mesg:
				if !ok {
					return
				}
				output <- item

				if res,ok:=item.Data.(*messages.JsonRpcResponse ); ok {					 
						s.mu.Lock()
						delete(s.pendingRequests, *res.GetId())
						empty := len(s.pendingRequests) == 0
						s.mu.Unlock()
						if empty {
							return
						} 
				} 
			}
		}
	}()
	return output 
}
