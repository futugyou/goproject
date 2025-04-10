package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type IMcpServerPrompt interface {
	IMcpServerPrimitive
	GetProtocolPrompt() *types.Prompt
	Get(ctx context.Context, request RequestContext[*types.GetPromptRequestParams]) (*types.GetPromptResult, error)
}
