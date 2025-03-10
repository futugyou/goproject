package functions

import (
	"context"

	"github.com/futugyou/ai-extension/abstractions"
)

type AIFunction interface {
	abstractions.AITool
	Invoke(ctx context.Context, arguments map[string]interface{}) (interface{}, error)
}

type BaseAIFunction struct {
	abstractions.BaseAITool
}

func (f *BaseAIFunction) Invoke(ctx context.Context, arguments map[string]interface{}) (interface{}, error) {
	return f.InvokeCore(ctx, arguments)
}

func (f *BaseAIFunction) InvokeCore(ctx context.Context, arguments map[string]interface{}) (interface{}, error) {
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

func (f *ExampleFunction) InvokeCore(arguments map[string]interface{}, ctx context.Context) (interface{}, error) {
	return "run ExampleFunction", nil
}
