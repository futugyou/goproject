package sse

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type ItemFormatter[T any] func(item SseItem[T], w *bufio.Writer)

func Write[T any](ctx context.Context, source <-chan SseItem[T], dst io.Writer, itemFormatter ItemFormatter[T]) error {
	writer := bufio.NewWriter(dst)

	if itemFormatter == nil {
		var zero T
		if reflect.TypeOf(zero).Kind() == reflect.String {
			anyFormatter := any(DefaultStringFormatter)
			itemFormatter = anyFormatter.(ItemFormatter[T])
		} else {
			return fmt.Errorf("no formatter provided for type %T", zero)
		}
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-source:
			if !ok {
				return writer.Flush()
			}

			var dataBuilder strings.Builder
			itemFormatter(item, bufio.NewWriter(&dataBuilder))

			err := writeSseEvent(writer, item, dataBuilder.String())
			if err != nil {
				return err
			}

			err = writer.Flush()
			if err != nil {
				return err
			}
		}
	}
}
func writeSseEvent[T any](w *bufio.Writer, item SseItem[T], data string) error {
	if item.EventType != "" {
		if strings.ContainsAny(item.EventType, "\r\n") {
			return fmt.Errorf("event type cannot contain line breaks")
		}
		_, _ = fmt.Fprintf(w, "event: %s\n", item.EventType)
	}

	writeLinesWithPrefix(w, "data: ", data)
	_, _ = fmt.Fprint(w, "\n")

	if item.EventId != nil && *item.EventId != "" {
		if strings.ContainsAny(*item.EventId, "\r\n") {
			return fmt.Errorf("event id cannot contain line breaks")
		}
		_, _ = fmt.Fprintf(w, "id: %s\n", *item.EventId)
	}

	if item.ReconnectionInterval > 0 {
		_, _ = fmt.Fprintf(w, "retry: %d\n", int(item.ReconnectionInterval.Milliseconds()))
	}

	_, _ = fmt.Fprint(w, "\n")
	return nil
}

func writeLinesWithPrefix(w *bufio.Writer, prefix string, data string) {
	lines := strings.Split(strings.ReplaceAll(data, "\r\n", "\n"), "\n")
	for _, line := range lines {
		_, _ = fmt.Fprintf(w, "%s%s\n", prefix, line)
	}
}

func DefaultStringFormatter(item SseItem[string], w *bufio.Writer) {
	_, _ = w.WriteString(item.Data)
}
