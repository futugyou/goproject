package httputils

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"strings"
)

const (
	DataPrefix  = "data: "
	LastToken   = "[DONE]"
	DoneMessage = DataPrefix + LastToken
)

// ParseStreamAsync parses the SSE data stream and returns the parsed data channel
func ParseStreamAsync[T any](ctx context.Context, stream io.Reader) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(stream)
		var buffer strings.Builder
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			line := scanner.Text()
			if strings.TrimSpace(line) == "" { // \n\n detected => Message delimiter
				if buffer.Len() == 0 {
					continue
				}
				message := buffer.String()
				buffer.Reset()
				if strings.TrimSpace(message) == DoneMessage {
					return
				}
				if value, ok := ParseMessage[T](message); ok {
					ch <- value
				}
			} else {
				buffer.WriteString(line + "\n")
			}
		}

		if buffer.Len() > 0 {
			message := buffer.String()
			if strings.TrimSpace(message) == DoneMessage {
				return
			}
			if value, ok := ParseMessage[T](message); ok {
				ch <- value
			}
		}
	}()
	return ch
}

// ParseMessage parses a single SSE message
func ParseMessage[T any](message string) (T, bool) {
	var result T
	if strings.TrimSpace(message) == "" {
		return result, false
	}

	var jsonBuilder strings.Builder
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, DataPrefix) {
			jsonBuilder.WriteString(line[len(DataPrefix):])
		}
	}

	if err := json.Unmarshal([]byte(jsonBuilder.String()), &result); err != nil {
		return result, false
	}
	return result, true
}
