package transport

import (
	"context"
	"net/http"

	"github.com/futugyou/yomawari/mcp/configuration"
)

type SseClientTransport struct {
	options        *SseClientTransportOptions
	serverConfig   *configuration.McpServerConfig
	httpClient     *http.Client
	ownsHttpClient bool
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
	return transport
}

// Connect implements IClientTransport.
func (s *SseClientTransport) Connect(ctx context.Context) (ITransport, error) {
	sessionTransport := NewSseClientSessionTransport(s.serverConfig, s.options, s.httpClient)
	err := sessionTransport.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return sessionTransport, nil
}
