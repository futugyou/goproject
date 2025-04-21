package sse

import "time"

type SseItemParser[T any] func(eventType string, data []byte) T

type SseItem[T any] struct {
	EventType            string
	EventId              *string
	ReconnectionInterval time.Duration
	Data                 T
}

func NewSseItem[T any](data T, eventType string) *SseItem[T] {
	s := &SseItem[T]{
		Data:      data,
		EventType: eventType,
	}
	if len(s.EventType) == 0 {
		s.EventType = "message"
	}
	return s
}
