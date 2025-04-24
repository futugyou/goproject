package shared

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/transport"
)

type NullProgress struct {
}

func (p *NullProgress) Report(value transport.ProgressNotificationValue) {}

type TokenProgress struct {
	endpoint      IMcpEndpoint
	progressToken transport.ProgressToken
}

func NewTokenProgress(endpoint IMcpEndpoint, progressToken transport.ProgressToken) *TokenProgress {
	return &TokenProgress{endpoint, progressToken}
}

func (p *TokenProgress) Report(value transport.ProgressNotificationValue) {
	p.endpoint.NotifyProgress(context.Background(), p.progressToken, value)
}
