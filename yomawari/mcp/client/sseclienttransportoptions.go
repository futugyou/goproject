package client

import (
	"net/url"
	"time"
)

type HttpTransportMode string

var HttpTransportModeAutoDetect HttpTransportMode = "AutoDetect"
var HttpTransportModeStreamableHttp HttpTransportMode = "StreamableHttp"
var HttpTransportModeSse HttpTransportMode = "Sse"

type SseClientTransportOptions struct {
	Endpoint             url.URL
	Name                 *string
	ConnectionTimeout    time.Duration
	AdditionalHeaders    map[string]string
	MaxReconnectAttempts int
	ReconnectDelay       time.Duration
	HttpTransportMode    HttpTransportMode
}

func NewSseClientTransportOptions() *SseClientTransportOptions {
	return &SseClientTransportOptions{
		ConnectionTimeout:    30 * time.Second,
		AdditionalHeaders:    map[string]string{},
		MaxReconnectAttempts: 3,
		ReconnectDelay:       5 * time.Second,
		HttpTransportMode:    HttpTransportModeAutoDetect,
	}
}
