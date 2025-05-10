package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type IMcpServerResource interface {
	IMcpServerPrimitive
	GetProtocolResourceTemplate() types.ResourceTemplate
	GetProtocolResource() *types.Resource
	Read(ctx context.Context, request RequestContext[*types.ReadResourceRequestParams]) (*types.ReadResourceResult, error)
}
