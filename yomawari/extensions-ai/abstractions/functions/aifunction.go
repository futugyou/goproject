package functions

import (
	"context"

	"github.com/futugyou/yomawari/extensions-ai/abstractions"
)

type AIFunction interface {
	abstractions.AITool
	Invoke(ctx context.Context, arguments AIFunctionArguments) (interface{}, error)
	GetJsonSchema() map[string]interface{}
}

type BaseAIFunction struct {
	abstractions.BaseAITool
	arguments AIFunctionArguments
}

func (t BaseAIFunction) GetParameters() map[string]interface{} {
	return map[string]interface{}{}
}

func (t BaseAIFunction) GetJsonSchema() map[string]interface{} {
	panic("not implemented")
}

func (f *BaseAIFunction) Invoke(ctx context.Context, arguments AIFunctionArguments) (interface{}, error) {
	return f.InvokeCore(ctx, arguments)
}

func (f *BaseAIFunction) InvokeCore(ctx context.Context, arguments AIFunctionArguments) (interface{}, error) {
	panic("InvokeCore must be implemented by subclass")
}

type ExampleFunction struct {
	BaseAIFunction
}

func NewExampleFunction() *ExampleFunction {
	return &ExampleFunction{
		BaseAIFunction{BaseAITool: abstractions.NewBaseAITool("ExampleFunction")},
	}
}

func (f *ExampleFunction) InvokeCore(ctx context.Context, arguments AIFunctionArguments) (interface{}, error) {
	return "run ExampleFunction", nil
}
