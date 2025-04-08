package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type CompletionsCapability struct {
	CompleteHandler func(context.Context, RequestContext[*types.CompleteRequestParams]) (*types.CompleteResult, error)
}
