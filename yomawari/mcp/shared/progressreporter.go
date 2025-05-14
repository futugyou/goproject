package shared

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type NullProgress struct {
}

func (p *NullProgress) Report(value transport.ProgressNotificationValue) {}

type TokenProgress struct {
	endpoint      IMcpEndpoint
	progressToken types.ProgressToken
}

func NewTokenProgress(endpoint IMcpEndpoint, progressToken types.ProgressToken) *TokenProgress {
	return &TokenProgress{endpoint, progressToken}
}

func (p *TokenProgress) Report(value transport.ProgressNotificationValue) {
	p.endpoint.NotifyProgress(context.Background(), p.progressToken, value)
}
