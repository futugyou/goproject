package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/futugyou/yomawari/extensions-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions-ai/abstractions/functions"
	"github.com/futugyou/yomawari/mcp"
	"github.com/futugyou/yomawari/mcp/protocol/types"
)

const RequestContextKey string = "__temporary_RequestContext"

type AIFunctionMcpServerPrompt struct {
	*McpServerPrompt
	AIFunction functions.AIFunction
}

func AIFunctionMcpServerPromptCreate(function functions.AIFunction, options McpServerPromptCreateOptions) *AIFunctionMcpServerPrompt {
	args := []types.PromptArgument{}

	if properties, ok := function.GetJsonSchema()["properties"].(json.RawMessage); ok {
		_ = properties
		// TODO: fill args
	}

	prompt := types.Prompt{
		Arguments:   args,
		Description: options.Description,
	}

	if options.Name != nil {
		prompt.Name = *options.Name
	}

	return &AIFunctionMcpServerPrompt{
		McpServerPrompt: &McpServerPrompt{
			ProtocolPrompt: prompt,
		},
		AIFunction: function,
	}
}

func (m *AIFunctionMcpServerPrompt) Get(ctx context.Context, request RequestContext[*types.GetPromptRequestParams]) (*types.GetPromptResult, error) {
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

	// TODO: Once we shift to the real AIFunctionFactory, the request should be passed via AIFunctionArguments.Context.
	arguments := map[string]interface{}{}
	if request.Params.Arguments != nil {
		arguments = request.Params.Arguments
	}
	arguments["RequestContextKey"] = request

	result, err := m.AIFunction.Invoke(ctx, arguments)
	if err != nil {
		return nil, err
	}
	switch r := result.(type) {
	case *string:
		return &types.GetPromptResult{
			Description: m.ProtocolPrompt.Description,
			Messages: []types.PromptMessage{{
				Content: types.Content{
					Type: "text",
					Text: r,
				},
				Role: types.RoleUser,
			}},
		}, nil
	case *types.GetPromptResult:
		return r, nil
	case *types.PromptMessage:
		return &types.GetPromptResult{
			Description: m.ProtocolPrompt.Description,
			Messages:    []types.PromptMessage{*r},
		}, nil
	case []types.PromptMessage:
		return &types.GetPromptResult{
			Description: m.ProtocolPrompt.Description,
			Messages:    r,
		}, nil
	case *chatcompletion.ChatMessage:
		return &types.GetPromptResult{
			Description: m.ProtocolPrompt.Description,
			Messages:    mcp.ChatMessageToPromptMessages(*r),
		}, nil
	case []chatcompletion.ChatMessage:
		msg := []types.PromptMessage{}
		for _, r := range r {
			msg = append(msg, mcp.ChatMessageToPromptMessages(r)...)
		}
		return &types.GetPromptResult{
			Description: m.ProtocolPrompt.Description,
			Messages:    msg,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
}
