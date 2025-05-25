package abstractions

import (
	"context"
)

type PromptRenderContext struct {
	Ctx               context.Context
	CancelFunc        context.CancelFunc
	IsStreaming       bool
	Kernel            Kernel
	Function          KernelFunction
	Arguments         KernelArguments
	Result            FunctionResult
	ExecutionSettings PromptExecutionSettings
	RenderedPrompt    string
}

type IPromptRenderFilter interface {
	OnPromptRender(context PromptRenderContext, next func(PromptRenderContext) error) error
}
