package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type IMcpServerResource interface {
	IMcpServerPrimitive
	GetProtocolResourceTemplate() protocol.ResourceTemplate
	GetProtocolResource() *protocol.Resource
	Read(ctx context.Context, request RequestContext[*protocol.ReadResourceRequestParams]) (*protocol.ReadResourceResult, error)
}
