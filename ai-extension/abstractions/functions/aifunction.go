package functions

import (
	"context"

	"github.com/futugyou/ai-extension/abstractions"
)

type AIFunction interface {
	abstractions.AITool
	InvokeAsync(ctx context.Context, arguments map[string]interface{}) (interface{}, error)
}

type BaseAIFunction struct {
	abstractions.BaseAITool
}

func (f *BaseAIFunction) InvokeAsync(ctx context.Context, arguments map[string]interface{}) (interface{}, error) {
	return f.InvokeCoreAsync(ctx, arguments)
}

func (f *BaseAIFunction) InvokeCoreAsync(ctx context.Context, arguments map[string]interface{}) (interface{}, error) {
	panic("InvokeCoreAsync must be implemented by subclass")
}

type ExampleFunction struct {
	BaseAIFunction
}

func NewExampleFunction() *ExampleFunction {
	return &ExampleFunction{
		BaseAIFunction{BaseAITool: abstractions.NewBaseAITool("ExampleFunction")},
	}
}

func (f *ExampleFunction) InvokeCoreAsync(arguments map[string]interface{}, ctx context.Context) (interface{}, error) {
	return "run ExampleFunction", nil
}
