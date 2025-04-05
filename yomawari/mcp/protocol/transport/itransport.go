package transport

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

type ITransport interface {
	IsConnected() bool
	MessageReader() <-chan messages.IJsonRpcMessage
	SendMessage(ctx context.Context, message messages.IJsonRpcMessage) error
	Close() error
}

const TransportTypesStdIo string = "stdio"
const TransportTypesSse string = "sse"

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

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				if data != "" {
					ch <- SseItem{EventType: eventType, Data: data}
					data = ""
				}
				continue
			}

			if strings.HasPrefix(line, "event: ") {
				eventType = strings.TrimPrefix(line, "event: ")
			} else if strings.HasPrefix(line, "data: ") {
				data += strings.TrimPrefix(line, "data: ") + "\n"
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return ch
}
