package client

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type McpClientResourceTemplate struct {
	client                   IMcpClient
	ProtocolResourceTemplate protocol.ResourceTemplate
	UriTemplate              string
	Name                     string
	Description              *string
	MimeTyp                  *string
}

func NewMcpClientResourceTemplate(client IMcpClient, protocolResourceTemplate protocol.ResourceTemplate) *McpClientResourceTemplate {
	return &McpClientResourceTemplate{
		client:                   client,
		ProtocolResourceTemplate: protocolResourceTemplate,
		UriTemplate:              protocolResourceTemplate.UriTemplate,
		Name:                     protocolResourceTemplate.Name,
		Description:              protocolResourceTemplate.Description,
		MimeTyp:                  protocolResourceTemplate.MimeType,
	}
}

func (m *McpClientResourceTemplate) Read(ctx context.Context, arguments map[string]interface{}) (*protocol.ReadResourceResult, error) {
	return m.client.ReadResourceWithUriAndArguments(ctx, m.UriTemplate, arguments)
}
