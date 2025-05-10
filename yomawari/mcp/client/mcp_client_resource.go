package client

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/types"
)

type McpClientResource struct {
	client           IMcpClient
	ProtocolResource types.Resource
	Uri              string
	Name             string
	Description      *string
	MimeTyp          *string
}

func NewMcpClientResource(client IMcpClient, protocolResource types.Resource) *McpClientResource {
	return &McpClientResource{
		client:           client,
		ProtocolResource: protocolResource,
		Uri:              protocolResource.Uri,
		Name:             protocolResource.Name,
		Description:      protocolResource.Description,
		MimeTyp:          protocolResource.MimeType,
	}
}

func (m *McpClientResource) Read(ctx context.Context) (*types.ReadResourceResult, error) {
	return m.client.ReadResource(ctx, m.Uri)
}
