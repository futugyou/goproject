package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/futugyou/yomawari/mcp/protocol"
	"github.com/futugyou/yomawari/runtime/sse"
)

type SseClientSessionTransport struct {
	*protocol.TransportBase
	httpClient            *http.Client
	Options               *SseClientTransportOptions
	SseEndpoint           *url.URL
	messageEndpoint       *url.URL
	ctx                   context.Context
	cancelFunc            context.CancelFunc
	connectionEstablished chan bool
	EndpointName          string
	disposed              bool

	receiveTaskCompleted chan struct{}
}

func NewSseClientSessionTransport(endpointName string, options *SseClientTransportOptions, httpClient *http.Client, messageChannel chan protocol.IJsonRpcMessage) *SseClientSessionTransport {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	if options == nil {
		options = &SseClientTransportOptions{}
	}
	ctx, cancel := context.WithCancel(context.Background())
	transport := &SseClientSessionTransport{
		TransportBase:         protocol.NewTransportBase(endpointName, messageChannel),
		httpClient:            httpClient,
		Options:               options,
		SseEndpoint:           &options.Endpoint,
		ctx:                   ctx,
		cancelFunc:            cancel,
		connectionEstablished: make(chan bool),
		EndpointName:          endpointName,
		receiveTaskCompleted:  make(chan struct{}),
	}
	return transport
}

func (t *SseClientSessionTransport) GetTransportKind() protocol.TransportKind {
	return protocol.TransportKindSse
}

func (t *SseClientSessionTransport) Connect(ctx context.Context) error {
	if t.IsConnected() {
		return fmt.Errorf("transport is already connected")
	}

	go t.receiveMessages(t.ctx)
	connectionTimeout := time.Duration(5) * time.Second
	if t.Options != nil && t.Options.ConnectionTimeout != 0 {
		connectionTimeout = t.Options.ConnectionTimeout * time.Second
	}

	select {
	case <-t.connectionEstablished:
	case <-time.After(connectionTimeout):
		return t.Close()
	}
	return nil
}

func (s *SseClientSessionTransport) HandleEndpointEvent(data string) error {
	if len(data) == 0 {
		return nil
	}

	if strings.HasPrefix(data, "http://") || strings.HasPrefix(data, "https://") {
		messageEndpoint, err := url.Parse(data)
		if err != nil {
			return err
		}
		s.messageEndpoint = messageEndpoint
	} else {
		endpointUri := fmt.Sprintf("%s/%s", strings.TrimRight(strings.TrimSuffix(s.SseEndpoint.String(), "/sse"), "/"), strings.TrimLeft(data, "/"))
		messageEndpoint, err := url.Parse(endpointUri)
		if err != nil {
			return err
		}
		s.messageEndpoint = messageEndpoint
	}

	// Set connected state
	s.SetConnected(true)
	select {
	case <-s.connectionEstablished:
	default:
		close(s.connectionEstablished)
	}
	return nil
}

func (s *SseClientSessionTransport) ProcessSseMessage(ctx context.Context, data string) error {
	if !s.IsConnected() {
		return nil
	}

	message, err := protocol.UnmarshalJsonRpcMessage([]byte(data))
	if err != nil {
		return err
	}

	s.WriteMessage(ctx, message)
	return nil
}

func (s *SseClientSessionTransport) receiveMessages(ctx context.Context) error {
	defer close(s.receiveTaskCompleted)
	defer func() {
		s.SetConnected(false)
	}()
	reconnectAttempts := 0

	for !s.IsConnected() && ctx.Err() == nil {
		err := s.connectAndProcessMessages(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			reconnectAttempts++
			if reconnectAttempts >= s.Options.MaxReconnectAttempts {
				return fmt.Errorf("exceeded reconnect limit: %w", err)
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(s.Options.ReconnectDelay):
				continue
			}
		}
	}

	return nil
}

func (s *SseClientSessionTransport) connectAndProcessMessages(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.SseEndpoint.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/event-stream")
	CopyAdditionalHeaders(req, s.Options.AdditionalHeaders, "")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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
			case "endpoint":
				if err := s.HandleEndpointEvent(event.Data); err != nil {
					return err
				}
			case "message":
				if err := s.ProcessSseMessage(ctx, event.Data); err != nil {
					return err
				}
			}
		}
	}
}

func (t *SseClientSessionTransport) Close() error {
	if t.disposed {
		return nil
	}
	t.disposed = true

	t.cancelFunc()

	select {
	case <-t.receiveTaskCompleted:
	case <-time.After(5 * time.Second):
	}

	t.SetConnected(false)
	return nil
}

func (t *SseClientSessionTransport) SendMessage(ctx context.Context, message protocol.IJsonRpcMessage) error {
	if t.messageEndpoint == nil {
		return fmt.Errorf("transport not connected")
	}

	data, err := protocol.MarshalJsonRpcMessage(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	var messageId = "(no id)"
	if msgWithId, ok := message.(protocol.IJsonRpcMessageWithId); ok {
		id := msgWithId.GetId()
		if id != nil {
			messageId = id.String()
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.messageEndpoint.String(), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	CopyAdditionalHeaders(req, t.Options.AdditionalHeaders, "")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	responseContent := string(body)

	if strings.EqualFold(responseContent, "accepted") {
		fmt.Printf("SSE Transport Post Accepted: %s, Message ID: %s\n", t.messageEndpoint.String(), messageId)
	} else {
		fmt.Printf("SSE Transport Post Not Accepted: %s, Message ID: %s, Response: %s\n", t.messageEndpoint.String(), messageId, responseContent)
		return fmt.Errorf("failed to send message")
	}

	return nil
}
