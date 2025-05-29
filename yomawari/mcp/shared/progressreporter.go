package shared

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type NullProgress struct {
}

func (p *NullProgress) Report(value protocol.ProgressNotificationValue) {}

type TokenProgress struct {
	endpoint      IMcpEndpoint
	progressToken protocol.ProgressToken
}

func NewTokenProgress(endpoint IMcpEndpoint, progressToken protocol.ProgressToken) *TokenProgress {
	return &TokenProgress{endpoint, progressToken}
}

func (p *TokenProgress) Report(value protocol.ProgressNotificationValue) {
	p.endpoint.NotifyProgress(context.Background(), p.progressToken, value)
}
