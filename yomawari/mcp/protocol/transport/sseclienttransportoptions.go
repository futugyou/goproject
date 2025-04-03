package transport

import "time"

type SseClientTransportOptions struct {
	ConnectionTimeout    time.Duration
	MaxReconnectAttempts int
	ReconnectDelay       time.Duration
	AdditionalHeaders    map[string]string
}

func NewSseClientTransportOptions() *SseClientTransportOptions {
	return &SseClientTransportOptions{
		ConnectionTimeout:    30 * time.Second,
		MaxReconnectAttempts: 3,
		ReconnectDelay:       5 * time.Second,
		AdditionalHeaders:    map[string]string{},
	}
}
