package mcp

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
)

type NullProgress struct {
}

func (p *NullProgress) Report(value messages.ProgressNotificationValue) {}

type TokenProgress struct {
	endpoint      IMcpEndpoint
	progressToken messages.ProgressToken
}

func NewTokenProgress(endpoint IMcpEndpoint, progressToken messages.ProgressToken) *TokenProgress {
	return &TokenProgress{endpoint, progressToken}
}

func (p *TokenProgress) Report(value messages.ProgressNotificationValue) {
	p.endpoint.NotifyProgress(context.Background(), p.progressToken, value)
}
