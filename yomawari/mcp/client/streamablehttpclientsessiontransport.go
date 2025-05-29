package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/futugyou/yomawari/mcp/protocol"
	"github.com/futugyou/yomawari/runtime/sse"
)

type StreamableHttpClientSessionTransport struct {
	*protocol.TransportBase
	httpClient   *http.Client
	Options      *SseClientTransportOptions
	ctx          context.Context
	cancelFunc   context.CancelFunc
	mcpSessionId string
	disposed     bool

	getReceiveTask chan struct{}
}

func NewStreamableHttpClientSessionTransport(httpClient *http.Client, options *SseClientTransportOptions, name string) *StreamableHttpClientSessionTransport {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if options == nil {
		options = &SseClientTransportOptions{}
	}
	ctx, cancel := context.WithCancel(context.Background())
	base := protocol.ClientTransportBase()
	base.Name = name
	transport := &StreamableHttpClientSessionTransport{
		TransportBase:  base,
		httpClient:     httpClient,
		Options:        options,
		ctx:            ctx,
		cancelFunc:     cancel,
		getReceiveTask: nil,
		disposed:       false,
	}
	transport.SetConnected(true)
	return transport
}

func (t *StreamableHttpClientSessionTransport) GetTransportKind() protocol.TransportKind {
	return protocol.TransportKindSse
}

func (t *StreamableHttpClientSessionTransport) SendMessage(ctx context.Context, message protocol.IJsonRpcMessage) error {
	var err error
	ctx, _ = protocol.MergeContexts(t.ctx, ctx)

	data, err := protocol.MarshalJsonRpcMessage(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.Options.Endpoint.String(), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json, text/event-stream")
	CopyAdditionalHeaders(req, t.Options.AdditionalHeaders, t.mcpSessionId)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rpcMessage protocol.IJsonRpcMessage
	var rpcRequest *protocol.JsonRpcRequest

	switch resp.Header.Get("Content-Type") {
	case "application/json":
		rpcMessage, err = t.processResponse(ctx, resp)
	case "text/event-stream":
		if condition, ok := message.(*protocol.JsonRpcRequest); ok && condition != nil {
			rpcRequest = condition
		}
		rpcMessage, err = t.processSseResponse(ctx, resp, rpcRequest)
	}

	if err != nil {
		return err
	}

	if rpcRequest == nil {
		return nil
	}

	if messageWithId, ok := rpcMessage.(protocol.IJsonRpcMessageWithId); !ok || messageWithId.GetId().String() != rpcRequest.GetId().String() {
		return fmt.Errorf("streamable HTTP POST response completed without a reply to request with ID: %s", rpcRequest.GetId().String())
	}

	if _, ok := rpcMessage.(*protocol.JsonRpcResponse); ok && rpcRequest.Method == protocol.RequestMethods_Initialize {
		t.mcpSessionId = resp.Header.Get("mcp-session-id")
		go t.receiveUnsolicitedMessages()
	}

	return nil
}

func (t *StreamableHttpClientSessionTransport) Close() error {
	if t.disposed {
		return nil
	}
	t.disposed = true

	t.cancelFunc()

	if t.getReceiveTask != nil {
		select {
		case <-t.getReceiveTask:
		case <-time.After(5 * time.Second):
		}
	}

	if t.Options.HttpTransportMode != HttpTransportModeAutoDetect || t.getReceiveTask != nil {
		t.SetConnected(false)
	}

	return nil
}

func CopyAdditionalHeaders(req *http.Request, additionalHeaders map[string]string, mcpSessionId string) {
	req.Header.Set("Content-Type", "application/json")

	if mcpSessionId != "" {
		req.Header.Set("mcp-session-id", mcpSessionId)
	}
	for k, v := range additionalHeaders {
		req.Header.Set(k, v)
	}
}

func (t *StreamableHttpClientSessionTransport) processResponse(ctx context.Context, resp *http.Response) (protocol.IJsonRpcMessage, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rpcMessage, err := protocol.UnmarshalJsonRpcMessage(body)
	if err != nil {
		return nil, err
	}

	err = t.WriteMessage(ctx, rpcMessage)
	if err != nil {
		return nil, err
	}
	return rpcMessage, nil
}

func (t *StreamableHttpClientSessionTransport) processSseResponse(ctx context.Context, resp *http.Response, rpcRequest *protocol.JsonRpcRequest) (protocol.IJsonRpcMessage, error) {
	sseWriter := sse.CreateSseParser(resp.Body)
	eventCh, errCh := sseWriter.EnumerateStream(ctx)

	rpcRequestId := ""
	if rpcRequest != nil {
		rpcRequestId = rpcRequest.GetId().String()
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errCh:
			return nil, err
		case event, ok := <-eventCh:
			if !ok {
				return nil, nil
			}

			switch event.EventType {
			case "message":
				rpcMessage, err := protocol.UnmarshalJsonRpcMessage([]byte(event.Data))
				if err != nil {
					return nil, err
				}

				err = t.WriteMessage(ctx, rpcMessage)
				if err != nil {
					return nil, err
				}

				if rpcMessageWithId, ok := rpcMessage.(protocol.IJsonRpcMessageWithId); ok && rpcMessageWithId != nil {
					if rpcMessageWithId.GetId().String() == rpcRequestId {
						return rpcMessageWithId, nil
					}
				}
			}
		}
	}
}

func (t *StreamableHttpClientSessionTransport) receiveUnsolicitedMessages() {
	t.getReceiveTask = make(chan struct{})
	defer close(t.getReceiveTask)
	req, err := http.NewRequestWithContext(t.ctx, http.MethodGet, t.Options.Endpoint.String(), nil)
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/json, text/event-stream")
	CopyAdditionalHeaders(req, t.Options.AdditionalHeaders, t.mcpSessionId)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	t.processSseResponse(t.ctx, resp, nil)
}
