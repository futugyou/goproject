package protocol

import (
	"context"
	"net/http"

	"github.com/futugyou/yomawari/mcp/configuration"
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
func (s *SseClientTransport) Connect(ctx context.Context) (ITransport, error) {
	if s.options.UseStreamableHttp {
		return NewStreamableHttpClientSessionTransport(s.httpClient, s.options, s.name), nil
	}

	sessionTransport := NewSseClientSessionTransport(s.serverConfig, s.options, s.httpClient)
	err := sessionTransport.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return sessionTransport, nil
}
