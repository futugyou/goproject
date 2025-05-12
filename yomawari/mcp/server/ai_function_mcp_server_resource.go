package server

import (
	"context"
	"fmt"
	"reflect"

	"github.com/futugyou/yomawari/extensions-ai/abstractions/contents"
	"github.com/futugyou/yomawari/extensions-ai/abstractions/functions"
	"github.com/futugyou/yomawari/mcp/protocol/types"
	"github.com/futugyou/yomawari/mcp/shared"
)

var _ IMcpServerResource = (*AIFunctionMcpServerResource)(nil)

type AIFunctionMcpServerResource struct {
	AIFunction       functions.AIFunction
	Resource         *types.Resource
	ResourceTemplate types.ResourceTemplate
	uriParser        *shared.UriParser
}

func NewAIFunctionMcpServerResource(function functions.AIFunction, resourceTemplate types.ResourceTemplate) *AIFunctionMcpServerResource {
	r := &AIFunctionMcpServerResource{
		AIFunction:       function,
		ResourceTemplate: resourceTemplate,
		Resource:         resourceTemplate.AsResource(),
	}
	r.uriParser, _ = shared.CreateUriParser(resourceTemplate.UriTemplate)
	return r
}

// GetId implements IMcpServerResource.
func (a *AIFunctionMcpServerResource) GetId() string {
	if a == nil || a.AIFunction == nil {
		return ""
	}

	return a.AIFunction.GetName()
}

// GetProtocolResource implements IMcpServerResource.
func (a *AIFunctionMcpServerResource) GetProtocolResource() *types.Resource {
	return a.Resource
}

// GetProtocolResourceTemplate implements IMcpServerResource.
func (a *AIFunctionMcpServerResource) GetProtocolResourceTemplate() types.ResourceTemplate {
	return a.ResourceTemplate
}

// Read implements IMcpServerResource.
func (m *AIFunctionMcpServerResource) Read(ctx context.Context, request RequestContext[*types.ReadResourceRequestParams]) (*types.ReadResourceResult, error) {
	if m == nil || m.AIFunction == nil {
		return nil, fmt.Errorf("ai function is nil")
	}

	if request.Params == nil {
		return nil, fmt.Errorf("request.Params is nil")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	arguments := functions.AIFunctionArguments{
		Context: map[interface{}]interface{}{reflect.TypeOf(request): request},
	}
	var matches map[string]string
	if m.uriParser != nil {
		matches = m.uriParser.Match(request.Params.Uri)
		for k, v := range matches {
			arguments.Set(k, v)
		}
	}

	result, err := m.AIFunction.Invoke(ctx, arguments)
	if err != nil {
		return nil, err
	}

	switch r := result.(type) {
	case *types.ReadResourceResult:
		return r, nil

	case types.IResourceContents:
		return &types.ReadResourceResult{
			Contents: []types.IResourceContents{r},
		}, nil
	case []types.IResourceContents:
		return &types.ReadResourceResult{
			Contents: r,
		}, nil

	case *contents.TextContent:
		return &types.ReadResourceResult{
			Contents: []types.IResourceContents{
				&types.TextResourceContents{
					BaseResourceContents: types.BaseResourceContents{
						Uri:      request.Params.Uri,
						MimeType: m.ResourceTemplate.MimeType,
					},
					Text: r.Text,
				},
			},
		}, nil

	case *contents.DataContent:
		return &types.ReadResourceResult{
			Contents: []types.IResourceContents{
				&types.BlobResourceContents{
					BaseResourceContents: types.BaseResourceContents{Uri: request.Params.Uri, MimeType: &r.MediaType},
					Blob:                 string(r.GetBase64Data()),
				},
			},
		}, nil

	case *string:
		return &types.ReadResourceResult{
			Contents: []types.IResourceContents{
				&types.TextResourceContents{
					BaseResourceContents: types.BaseResourceContents{
						Uri:      request.Params.Uri,
						MimeType: m.ResourceTemplate.MimeType,
					},
					Text: *r,
				},
			},
		}, nil
	case []string:
		contents := []types.IResourceContents{}
		for _, v := range r {
			contents = append(contents, &types.TextResourceContents{
				BaseResourceContents: types.BaseResourceContents{
					Uri:      request.Params.Uri,
					MimeType: m.ResourceTemplate.MimeType,
				},
				Text: v,
			})
		}
		return &types.ReadResourceResult{
			Contents: contents,
		}, nil
	case []contents.IAIContent:
		conts := []types.IResourceContents{}
		for _, v := range r {
			if a, ok := v.(*contents.TextContent); ok {
				conts = append(conts, &types.TextResourceContents{
					BaseResourceContents: types.BaseResourceContents{
						Uri:      request.Params.Uri,
						MimeType: m.ResourceTemplate.MimeType,
					},
					Text: a.Text,
				})
			}
			if a, ok := v.(*contents.DataContent); ok {
				conts = append(conts,
					&types.BlobResourceContents{
						BaseResourceContents: types.BaseResourceContents{Uri: request.Params.Uri, MimeType: &a.MediaType},
						Blob:                 string(a.GetBase64Data()),
					})
			}
		}
		return &types.ReadResourceResult{
			Contents: conts,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
}
