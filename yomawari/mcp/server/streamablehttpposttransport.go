package server

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/futugyou/yomawari/mcp/protocol"
	"github.com/futugyou/yomawari/mcp/shared"
	"github.com/futugyou/yomawari/runtime/sse"
)

var _ protocol.ITransport = (*StreamableHttpPostTransport)(nil)

type StreamableHttpPostTransport struct {
	httpBodies      *shared.DuplexPipe
	sseWriter       *shared.SseWriter
	pendingRequests *protocol.RequestId
	parentTransport *StreamableHttpServerTransport
	mu              sync.Mutex
}

// GetTransportKind implements ITransport.
func (s *StreamableHttpPostTransport) GetTransportKind() protocol.TransportKind {
	return protocol.TransportKindHttp
}

func NewStreamableHttpPostTransport(parentTransport *StreamableHttpServerTransport, httpBodies *shared.DuplexPipe) *StreamableHttpPostTransport {

	return &StreamableHttpPostTransport{
		httpBodies:      httpBodies,
		sseWriter:       shared.NewSseWriter(""),
		parentTransport: parentTransport,
	}
}

// Close implements ITransport.
func (s *StreamableHttpPostTransport) Close() error {
	s.sseWriter.Dispose()
	return nil
}

// MessageReader implements ITransport.
func (s *StreamableHttpPostTransport) MessageReader() <-chan protocol.IJsonRpcMessage {
	panic("JsonRpcMessage.RelatedTransport should only be used for sending messages")
}

// SendMessage implements ITransport.
func (s *StreamableHttpPostTransport) SendMessage(ctx context.Context, message protocol.IJsonRpcMessage) error {
	if _, ok := message.(*protocol.JsonRpcRequest); ok && s.parentTransport.Stateless {
		return fmt.Errorf("server to client requests are not supported in stateless mode")
	}

	return s.sseWriter.SendMessage(ctx, message)
}

func (s *StreamableHttpPostTransport) Run(ctx context.Context) (bool, error) {
	data, err := io.ReadAll(s.httpBodies.Input)
	if err != nil {
		return false, err
	}
	msg, err := protocol.UnmarshalJsonRpcMessage(data)
	if err != nil {
		return false, err
	}
	if err := s.onMessageReceived(ctx, msg); err != nil {
		return false, err
	}

	s.mu.Lock()
	noRequests := s.pendingRequests == nil || s.pendingRequests.IsDefault()
	s.mu.Unlock()

	if noRequests {
		return false, nil
	}

	s.sseWriter.MessageFilter = s.stopOnFinalResponseFilter
	s.sseWriter.WriteAll(ctx, s.httpBodies.Output)
	return true, nil
}

//lint:ignore U1000 used by reflection or reserved for future use
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
			msg, err := protocol.UnmarshalJsonRpcMessage([]byte(raw))
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
		msg, err := protocol.UnmarshalJsonRpcMessage(data)
		if err != nil {
			return err
		}
		if err := s.onMessageReceived(ctx, msg); err != nil {
			return err
		}
	}

	return nil
}

func (s *StreamableHttpPostTransport) onMessageReceived(ctx context.Context, msg protocol.IJsonRpcMessage) error {
	if msg == nil {
		return fmt.Errorf("received invalid null message")
	}

	if request, ok := msg.(*protocol.JsonRpcRequest); ok {
		s.mu.Lock()
		s.pendingRequests = request.GetId()
		if s.parentTransport != nil && s.parentTransport.Stateless && request.Method == protocol.RequestMethods_Initialize && request.Params != nil {
			var r protocol.InitializeRequestParams
			var initialized bool

			if c, ok := request.Params.(*protocol.InitializeRequestParams); ok {
				r = *c
				initialized = true
			} else {
				data, err := json.Marshal(request.Params)
				if err == nil {
					if err := json.Unmarshal(data, &r); err == nil {
						initialized = true
					}
				}
			}

			if initialized {
				s.parentTransport.InitializeRequest = r
			}
		}
		s.mu.Unlock()
	}

	msg.SetRelatedTransport(s)

	if s.parentTransport == nil || s.parentTransport.incomingChannel == nil {
		return fmt.Errorf("incoming channel is nil")
	}

	select {
	case s.parentTransport.incomingChannel <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}

func (s *StreamableHttpPostTransport) stopOnFinalResponseFilter(ctx context.Context, mesg chan sse.SseItem[protocol.IJsonRpcMessage]) chan sse.SseItem[protocol.IJsonRpcMessage] {
	output := make(chan sse.SseItem[protocol.IJsonRpcMessage])
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

				if res, ok := item.Data.(*protocol.JsonRpcResponse); ok && res.Id == s.pendingRequests {
					return
				}
			}
		}
	}()
	return output
}
