package client

import (
	"net/url"
	"time"
)

type SseClientTransportOptions struct {
	Endpoint             url.URL
	Name                 *string
	ConnectionTimeout    time.Duration
	AdditionalHeaders    map[string]string
	MaxReconnectAttempts int
	ReconnectDelay       time.Duration
	UseStreamableHttp    bool
}

func NewSseClientTransportOptions() *SseClientTransportOptions {
	return &SseClientTransportOptions{
		ConnectionTimeout:    30 * time.Second,
		AdditionalHeaders:    map[string]string{},
		MaxReconnectAttempts: 3,
		ReconnectDelay:       5 * time.Second,
	}
}
