package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type IMcpServerPrompt interface {
	IMcpServerPrimitive
	GetProtocolPrompt() *protocol.Prompt
	Get(ctx context.Context, request RequestContext[*protocol.GetPromptRequestParams]) (*protocol.GetPromptResult, error)
}
