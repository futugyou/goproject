package transport

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"sync"

	"github.com/futugyou/yomawari/runtime/sse"
)

type SseWriter struct {
	messages        chan sse.SseItem[IJsonRpcMessage]
	mu              sync.Mutex
	ctx             context.Context
	cancelFunc      context.CancelFunc
	messageEndpoint string
	task            chan error
	disposed        bool
	MessageFilter   func(ctx context.Context, mesg chan sse.SseItem[IJsonRpcMessage]) chan sse.SseItem[IJsonRpcMessage]
}

func NewSseWriter(messageEndpoint string) *SseWriter {
	return &SseWriter{
		messages:        make(chan sse.SseItem[IJsonRpcMessage]),
		messageEndpoint: messageEndpoint,
	}
}

func (s *SseWriter) WriteAll(ctx context.Context, sseResponseStream io.Writer) chan error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.task != nil {
		select {
		case <-s.task:
		default:
			s.task <- fmt.Errorf("WriteAll already called")
			return s.task
		}
	}
	s.task = make(chan error, 1)

	if len(s.messageEndpoint) > 0 {
		select {
		case s.messages <- *sse.NewSseItem[IJsonRpcMessage](nil, "endpoint"):
		default:
			s.task <- fmt.Errorf("you must call RunAsync before calling SendMessage")
			return s.task
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	s.ctx = ctx
	s.cancelFunc = cancel

	var msg = s.messages
	if s.MessageFilter != nil {
		msg = s.MessageFilter(ctx, msg)
	}

	go func() {
		s.task <- sse.Write(ctx, msg, sseResponseStream, s.writeJsonRpcMessageToBuffer())
	}()

	return s.task
}

func (s *SseWriter) writeJsonRpcMessageToBuffer() sse.ItemFormatter[IJsonRpcMessage] {
	return func(item sse.SseItem[IJsonRpcMessage], writer *bufio.Writer) error {
		if item.EventType == "endpoint" && len(s.messageEndpoint) > 0 {
			_, err := fmt.Fprintf(writer, "%s", base64.URLEncoding.EncodeToString([]byte(s.messageEndpoint)))
			return err
		}
		d, err := MarshalJsonRpcMessage(item.Data)
		if err != nil {
			return err
		}
		_, err = writer.Write(d)
		return err
	}
}

func (s *SseWriter) SendMessage(ctx context.Context, message IJsonRpcMessage) error {
	if message == nil {
		return fmt.Errorf("message is nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.disposed {
		return fmt.Errorf("sseWriter is disposed")
	}

	var err error
	select {
	case s.messages <- *sse.NewSseItem[IJsonRpcMessage](message, "message"):
	default:
		err = fmt.Errorf("something went wrong sending the message")
	}

	return err
}

func (s *SseWriter) Dispose() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.disposed {
		return
	}

	if s.cancelFunc != nil {
		s.cancelFunc()
		s.cancelFunc = nil
	}

	close(s.messages)

	if s.task != nil {
		select {
		case <-s.task:
		default:
			s.task <- nil
		}
		close(s.task)
		s.task = nil
	}

	s.disposed = true
}
