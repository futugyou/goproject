package server

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
	"github.com/futugyou/yomawari/mcp"
	"github.com/futugyou/yomawari/mcp/protocol"
)

var _ IMcpServerPrompt = (*AIFunctionMcpServerPrompt)(nil)

type AIFunctionMcpServerPrompt struct {
	AIFunction     functions.AIFunction
	ProtocolPrompt *protocol.Prompt
}

// GetId implements IMcpServerPrompt.
func (m *AIFunctionMcpServerPrompt) GetId() string {
	if m == nil || m.AIFunction == nil {
		return ""
	}

	return m.AIFunction.GetName()
}

// GetProtocolPrompt implements IMcpServerPrompt.
func (m *AIFunctionMcpServerPrompt) GetProtocolPrompt() *protocol.Prompt {
	if m == nil {
		return nil
	}

	return m.ProtocolPrompt
}

func PromptCreate(function functions.AIFunction, options McpServerPromptCreateOptions) AIFunctionMcpServerPrompt {
	args := []protocol.PromptArgument{}
	requiredList := []string{}
	if requireds, ok := function.GetJsonSchema()["required"].(json.RawMessage); ok {
		json.Unmarshal(requireds, &requiredList)
	}

	if properties, ok := function.GetJsonSchema()["properties"].(json.RawMessage); ok {
		var propertyList []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := json.Unmarshal(properties, &propertyList); err == nil {
			for _, property := range propertyList {
				required := false
				for _, v := range requiredList {
					if v == property.Name {
						required = true
						break
					}
				}
				args = append(args, protocol.PromptArgument{
					Name:        property.Name,
					Description: &property.Description,
					Required:    &required,
				})
			}
		}
	}

	prompt := &protocol.Prompt{
		Arguments:   args,
		Description: options.Description,
	}

	if options.Name != nil {
		prompt.Name = *options.Name
	}

	return AIFunctionMcpServerPrompt{
		ProtocolPrompt: prompt,
		AIFunction:     function,
	}
}

func PromptDynamicCreate(methodInfo reflect.Value, options McpServerPromptCreateOptions) AIFunctionMcpServerPrompt {
	factory := functions.NewAIFunctionFactory()
	op := createPromptCreateOption(options, methodInfo)
	function, err := factory.Create(methodInfo, op)
	if err != nil {
		return AIFunctionMcpServerPrompt{}
	}
	return PromptCreate(function, options)
}

func createPromptCreateOption(options McpServerPromptCreateOptions, methodInfo reflect.Value) *functions.AIFunctionFactoryOptions {
	op := &functions.AIFunctionFactoryOptions{
		SerializerOptions:    &json.Encoder{},
		ParameterNames:       []string{},
		JSONSchemaOptions:    map[string]interface{}{},
		AdditionalProperties: map[string]interface{}{},
	}
	if options.Name != nil {
		op.Name = *options.Name
	} else {
		op.Name = methodInfo.Type().Method(0).Name
	}
	if options.Description != nil {
		op.Description = *options.Description
	}
	return op
}

func (m *AIFunctionMcpServerPrompt) Get(ctx context.Context, request RequestContext[*protocol.GetPromptRequestParams]) (*protocol.GetPromptResult, error) {
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
	case *string:
		return &protocol.GetPromptResult{
			Description: m.ProtocolPrompt.Description,
			Messages: []protocol.PromptMessage{{
				Content: protocol.Content{
					Type: "text",
					Text: r,
				},
				Role: protocol.RoleUser,
			}},
		}, nil
	case *protocol.GetPromptResult:
		return r, nil
	case *protocol.PromptMessage:
		return &protocol.GetPromptResult{
			Description: m.ProtocolPrompt.Description,
			Messages:    []protocol.PromptMessage{*r},
		}, nil
	case []protocol.PromptMessage:
		return &protocol.GetPromptResult{
			Description: m.ProtocolPrompt.Description,
			Messages:    r,
		}, nil
	case *chatcompletion.ChatMessage:
		return &protocol.GetPromptResult{
			Description: m.ProtocolPrompt.Description,
			Messages:    mcp.ChatMessageToPromptMessages(*r),
		}, nil
	case []chatcompletion.ChatMessage:
		msg := []protocol.PromptMessage{}
		for _, r := range r {
			msg = append(msg, mcp.ChatMessageToPromptMessages(r)...)
		}
		return &protocol.GetPromptResult{
			Description: m.ProtocolPrompt.Description,
			Messages:    msg,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}
}
