package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/futugyou/yomawari/mcp/configuration"
	"github.com/futugyou/yomawari/mcp/protocol"
)

var _ IClientTransport = (*SseClientTransport)(nil)

type SseClientTransport struct {
	options        *SseClientTransportOptions
	serverConfig   *configuration.McpServerConfig
	httpClient     *http.Client
	ownsHttpClient bool
	name           string
}

// GetName implements IClientTransport.
func (s *SseClientTransport) GetName() string {
	return s.name
}

func NewSseClientTransport(serverConfig *configuration.McpServerConfig, options *SseClientTransportOptions, httpClient *http.Client, ownsHttpClient *bool) *SseClientTransport {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	if serverConfig == nil {
		serverConfig = &configuration.McpServerConfig{}
	}
	if options == nil {
		options = &SseClientTransportOptions{}
	}
	transport := &SseClientTransport{
		options:        options,
		serverConfig:   serverConfig,
		httpClient:     httpClient,
		ownsHttpClient: false,
	}

	if ownsHttpClient != nil {
		transport.ownsHttpClient = *ownsHttpClient
	}

	if options.Name != nil {
		transport.name = *options.Name
	}

	if len(transport.name) == 0 {
		transport.name = options.Endpoint.String()
	}

	return transport
}

// Connect implements IClientTransport.
func (s *SseClientTransport) Connect(ctx context.Context) (protocol.ITransport, error) {

	switch s.options.HttpTransportMode {
	case HttpTransportModeAutoDetect:
		return NewAutoDetectingClientSessionTransport(s.httpClient, s.options, s.name), nil
	case HttpTransportModeStreamableHttp:
		return NewStreamableHttpClientSessionTransport(s.httpClient, s.options, s.name), nil
	case HttpTransportModeSse:
		return s.connectSseTransport(ctx)
	default:
		return nil, fmt.Errorf("unsupported transport mode: %s", s.options.HttpTransportMode)
	}
}

func (s *SseClientTransport) connectSseTransport(ctx context.Context) (protocol.ITransport, error) {
	sessionTransport := NewSseClientSessionTransport(s.name, s.options, s.httpClient, nil)
	err := sessionTransport.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return sessionTransport, nil
}
