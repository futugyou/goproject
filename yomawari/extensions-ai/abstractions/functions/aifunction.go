package functions

import (
	"context"

	"github.com/futugyou/yomawari/extensions-ai/abstractions"
)

var _ AIFunction = (*BaseAIFunction)(nil)

type AIFunction interface {
	abstractions.AITool
	Invoke(ctx context.Context, arguments AIFunctionArguments) (interface{}, error)
	// i want use https://github.com/invopop/jsonschema
	GetJsonSchema() map[string]interface{}
}

type BaseAIFunction struct {
	*abstractions.BaseAITool
}

func (t *BaseAIFunction) GetJsonSchema() map[string]interface{} {
	panic("not implemented")
}

func (f *BaseAIFunction) Invoke(ctx context.Context, arguments AIFunctionArguments) (interface{}, error) {
	return f.InvokeCore(ctx, arguments)
}

func (f *BaseAIFunction) InvokeCore(ctx context.Context, arguments AIFunctionArguments) (interface{}, error) {
	panic("InvokeCore must be implemented by subclass")
}
