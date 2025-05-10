package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

var _ IMcpServerResource = (*AIFunctionMcpServerResource)(nil)

type AIFunctionMcpServerResource struct {
}

// GetId implements IMcpServerResource.
func (a *AIFunctionMcpServerResource) GetId() string {
	panic("unimplemented")
}

// GetProtocolResource implements IMcpServerResource.
func (a *AIFunctionMcpServerResource) GetProtocolResource() *types.Resource {
	panic("unimplemented")
}

// GetProtocolResourceTemplate implements IMcpServerResource.
func (a *AIFunctionMcpServerResource) GetProtocolResourceTemplate() types.ResourceTemplate {
	panic("unimplemented")
}

// Read implements IMcpServerResource.
func (a *AIFunctionMcpServerResource) Read(ctx context.Context, request RequestContext[*types.ReadResourceRequestParams]) (*types.ReadResourceResult, error) {
	panic("unimplemented")
}
