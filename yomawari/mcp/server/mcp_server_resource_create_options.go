package server

import "github.com/futugyou/yomawari/core"

type McpServerResourceCreateOptions struct {
	Services    core.IServiceProvider
	UriTemplate *string
	Name        *string
	Description *string
	MimeType    *string
}

func (m McpServerResourceCreateOptions) Clone() McpServerResourceCreateOptions {
	return McpServerResourceCreateOptions{
		Services:    m.Services,
		UriTemplate: m.UriTemplate,
		Name:        m.Name,
		Description: m.Description,
		MimeType:    m.MimeType,
	}
}
