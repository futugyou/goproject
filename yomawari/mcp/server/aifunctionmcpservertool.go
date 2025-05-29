package server

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
	"github.com/futugyou/yomawari/mcp"
	"github.com/futugyou/yomawari/mcp/protocol"
)

var _ IMcpServerTool = (*AIFunctionMcpServerTool)(nil)

type AIFunctionMcpServerTool struct {
	ProtocolTool *protocol.Tool
	AIFunction   functions.AIFunction
}

// GetId implements IMcpServerTool.
func (a *AIFunctionMcpServerTool) GetId() string {
	if a == nil || a.ProtocolTool == nil {
		return ""
	}
	return a.ProtocolTool.Name
}

// GetProtocolTool implements IMcpServerTool.
func (a *AIFunctionMcpServerTool) GetProtocolTool() *protocol.Tool {
	if a == nil {
		return nil
	}
	return a.ProtocolTool
}

// Invoke implements IMcpServerTool.
func (m *AIFunctionMcpServerTool) Invoke(ctx context.Context, request RequestContext[*protocol.CallToolRequestParams]) (*protocol.CallToolResult, error) {

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

	if request.Params.Arguments != nil {
		for key, v := range request.Params.Arguments {
			arguments.Set(key, v)
		}
	}

	result, err := m.AIFunction.Invoke(ctx, arguments)
	if err != nil {
		return nil, err
	}
	switch r := result.(type) {
	case contents.IAIContent:
		isError := false
		if _, ok := r.(contents.ErrorContent); ok {
			isError = true
		}
		rr := protocol.NewCallToolResultWithContent(mcp.AIContentToContent(r))
		rr.IsError = isError
		return rr, nil
	case *string:
		result := protocol.NewCallToolResult()
		result.Content = append(result.Content, protocol.Content{Type: "text", Text: r})
		return result, nil
	case *protocol.Content:
		return protocol.NewCallToolResultWithContent(*r), nil
	case []string:
		result := protocol.NewCallToolResult()
		for _, v := range r {
			result.Content = append(result.Content, protocol.Content{Type: "text", Text: &v})
		}
		return result, nil
	case []contents.IAIContent:
		return onvertAIContentEnumerableToCallToolResponse(r), nil
	case []protocol.Content:
		return protocol.NewCallToolResultWithContents(r), nil
	case *protocol.CallToolResult:
		return r, nil
	default:
		// how to marshal?
		data, err := json.Marshal(r)
		if err != nil {
			return nil, err
		}
		text := string(data)
		result := protocol.NewCallToolResult()
		result.Content = append(result.Content, protocol.Content{Type: "text", Text: &text})
		return result, nil
	}
}

func onvertAIContentEnumerableToCallToolResponse(contentItems []contents.IAIContent) *protocol.CallToolResult {
	contentList := []protocol.Content{}
	allErrorContent := true
	hasAny := false

	for _, item := range contentItems {
		contentList = append(contentList, mcp.AIContentToContent(item))
		hasAny = true

		if _, ok := item.(contents.ErrorContent); !ok && allErrorContent {
			allErrorContent = false
		}
	}

	return &protocol.CallToolResult{
		Content: contentList,
		IsError: allErrorContent && hasAny,
	}
}
