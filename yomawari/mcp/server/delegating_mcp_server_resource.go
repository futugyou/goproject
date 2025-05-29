package server

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

var _ IMcpServerResource = (*DelegatingMcpServerResource)(nil)

type DelegatingMcpServerResource struct {
	delegate IMcpServerResource
}

func NewDelegatingMcpServerResource(delegate IMcpServerResource) *DelegatingMcpServerResource {
	return &DelegatingMcpServerResource{
		delegate: delegate,
	}
}

// GetId implements IMcpServerResource.
func (d *DelegatingMcpServerResource) GetId() string {
	return d.delegate.GetId()
}

// GetProtocolResource implements IMcpServerResource.
func (d *DelegatingMcpServerResource) GetProtocolResource() *protocol.Resource {
	return d.delegate.GetProtocolResource()
}

// GetProtocolResourceTemplate implements IMcpServerResource.
func (d *DelegatingMcpServerResource) GetProtocolResourceTemplate() protocol.ResourceTemplate {
	return d.delegate.GetProtocolResourceTemplate()
}

// Read implements IMcpServerResource.
func (d *DelegatingMcpServerResource) Read(ctx context.Context, request RequestContext[*protocol.ReadResourceRequestParams]) (*protocol.ReadResourceResult, error) {
	return d.delegate.Read(ctx, request)
}
