package client

import (
	"context"
	"encoding/json"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
	"github.com/futugyou/yomawari/mcp/protocol"
)

var _ functions.AIFunction = &McpClientTool{}

type McpClientTool struct {
	functions.BaseAIFunction
	additionalProperties map[string]interface{}
	client               IMcpClient
	name                 string
	description          *string
	ProtocolTool         protocol.Tool
}

func NewMcpClientTool(client IMcpClient, name string, description *string, protocolTool protocol.Tool) *McpClientTool {
	return &McpClientTool{
		additionalProperties: map[string]interface{}{"Strict": false},
		client:               client,
		name:                 name,
		description:          description,
		ProtocolTool:         protocolTool,
	}
}

func (m *McpClientTool) WithName(name string) *McpClientTool {
	return NewMcpClientTool(m.client, name, m.description, m.ProtocolTool)
}

func (m *McpClientTool) WithDescription(description string) *McpClientTool {
	return NewMcpClientTool(m.client, m.name, &description, m.ProtocolTool)
}

// GetAdditionalProperties implements functions.AIFunction.
func (m *McpClientTool) GetAdditionalProperties() map[string]interface{} {
	return m.additionalProperties
}

// GetDescription implements functions.AIFunction.
func (m *McpClientTool) GetDescription() string {
	if m == nil || m.description == nil {
		return ""
	}
	return *m.description
}

// GetJsonSchema implements functions.AIFunction.
func (m *McpClientTool) GetJsonSchema() map[string]interface{} {
	var result map[string]interface{}
	if err := json.Unmarshal(m.ProtocolTool.InputSchema, &result); err != nil {
		return nil
	}

	return result
}

// GetName implements functions.AIFunction.
func (m *McpClientTool) GetName() string {
	if len(m.name) == 0 {
		return "McpClientTool"
	}
	return m.name
}

// Invoke implements functions.AIFunction.
func (m *McpClientTool) Invoke(ctx context.Context, arguments functions.AIFunctionArguments) (interface{}, error) {
	return m.client.CallTool(ctx, m.ProtocolTool.Name, arguments.Items(), nil)
}
