package filters

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/ai_functional"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/functions"
)

type PromptRenderContext struct {
	Ctx               context.Context
	CancelFunc        context.CancelFunc
	IsStreaming       bool
	Kernel            abstractions.Kernel
	Function          functions.KernelFunction
	Arguments         functions.KernelArguments
	Result            functions.FunctionResult
	ExecutionSettings ai_functional.PromptExecutionSettings
	RenderedPrompt    string
}

type IPromptRenderFilter interface {
	OnPromptRender(context PromptRenderContext, next func(PromptRenderContext) error) error
}
