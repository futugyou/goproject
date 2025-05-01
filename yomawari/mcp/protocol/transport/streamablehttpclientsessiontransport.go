package transport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/futugyou/yomawari/runtime/sse"
)

type StreamableHttpClientSessionTransport struct {
	TransportBase
	httpClient   *http.Client
	Options      *SseClientTransportOptions
	ctx          context.Context
	cancelFunc   context.CancelFunc
	mcpSessionId string

	getReceiveTask chan struct{}
}

func NewStreamableHttpClientSessionTransport(httpClient *http.Client, options *SseClientTransportOptions) *StreamableHttpClientSessionTransport {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if options == nil {
		options = &SseClientTransportOptions{}
	}
	ctx, cancel := context.WithCancel(context.Background())
	transport := &StreamableHttpClientSessionTransport{
		TransportBase: TransportBase{
			messageChannel: make(chan IJsonRpcMessage),
			isConnected:    false,
		},
		httpClient:     httpClient,
		Options:        options,
		ctx:            ctx,
		cancelFunc:     cancel,
		getReceiveTask: make(chan struct{}),
	}
	transport.SetConnected(true)
	return transport
}

func (t *StreamableHttpClientSessionTransport) SendMessage(ctx context.Context, message IJsonRpcMessage) error {
	ctx, _ = mergeContexts(t.ctx, ctx)

	data, err := MarshalJsonRpcMessage(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.Options.Endpoint.String(), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	if t.mcpSessionId != "" {
		req.Header.Set("mcp-session-id", t.mcpSessionId)
	}
	for k, v := range t.Options.AdditionalHeaders {
		req.Header.Set(k, v)
	}
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rpcMessage IJsonRpcMessage
	var rpcRequest *JsonRpcRequest
	switch resp.Header.Get("Content-Type") {
	case "application/json":
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if rpcMessage, err = UnmarshalJsonRpcMessage(body); err != nil {
			return fmt.Errorf("failed to initialize client")
		}

		t.WriteMessage(ctx, rpcMessage)
	case "text/event-stream":
		if condition, ok := message.(*JsonRpcRequest); ok && condition != nil {
			rpcRequest = condition
		}

		sseWriter := sse.CreateSseParser(resp.Body)
		eventCh, errCh := sseWriter.EnumerateStream(ctx)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errCh:
				return err
			case event, ok := <-eventCh:
				if !ok {
					return nil
				}

				switch event.EventType {
				case "message":
					if rpcMessage, err = UnmarshalJsonRpcMessage([]byte(event.Data)); err != nil {
						return fmt.Errorf("failed to initialize client")
					}
					if rpcMessageWithId, ok := rpcMessage.(IJsonRpcMessageWithId); ok && rpcRequest != nil && rpcMessageWithId != nil {
						if rpcMessageWithId.GetId().String() == rpcRequest.GetId().String() {
							rpcMessage = rpcMessageWithId
						}
					}
				}
			}
		}
	}

	if rpcRequest == nil {
		return nil
	}

	if messageWithId, ok := rpcMessage.(IJsonRpcMessageWithId); !ok || messageWithId.GetId().String() != rpcRequest.GetId().String() {
		return fmt.Errorf("streamable HTTP POST response completed without a reply to request with ID: %s", rpcRequest.GetId().String())
	}

	if _, ok := rpcMessage.(*JsonRpcResponse); ok && rpcRequest.Method == RequestMethods_Initialize {
		t.mcpSessionId = resp.Header.Get("mcp-session-id")
	}
	return nil
}
