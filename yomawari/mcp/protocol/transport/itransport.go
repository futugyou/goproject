package transport

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

type ITransport interface {
	MessageReader() <-chan messages.IJsonRpcMessage
	SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error
	Close() error
}

const TransportTypesStdIo string = "stdio"
const TransportTypesSse string = "sse"

// TODO: use github.com/futugyou/yomawari/runtime/sse instead
type SseItem struct {
	EventType string
	Data      string
}

func ParseSSEStream(ctx context.Context, reader io.Reader) <-chan SseItem {
	ch := make(chan SseItem)

	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(reader)
		var eventType, data string

		lines := make(chan string)

		go func() {
			defer close(lines)
			for scanner.Scan() {
				select {
				case lines <- scanner.Text():
				case <-ctx.Done():
					return
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-lines:
				if !ok {
					return
				}

				if line == "" {
					if data != "" {
						ch <- SseItem{EventType: eventType, Data: strings.TrimSuffix(data, "\n")}
						data = ""
					}
					continue
				}

				if strings.HasPrefix(line, "event: ") {
					eventType = strings.TrimPrefix(line, "event: ")
				} else if strings.HasPrefix(line, "data: ") {
					data += strings.TrimPrefix(line, "data: ") + "\n"
				}
			}
		}
	}()

	return ch
}
