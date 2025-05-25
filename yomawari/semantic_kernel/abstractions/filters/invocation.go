package filters

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/ai_functional"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/functions"
)

type AutoFunctionInvocationContext struct {
	Ctx                   context.Context
	CancelFunc            context.CancelFunc
	IsStreaming           bool
	Arguments             functions.KernelArguments
	RequestSequenceIndex  int
	FunctionSequenceIndex int
	FunctionCount         int
	ToolCallId            string
	ChatMessageContent    contents.ChatMessageContent
	ExecutionSettings     ai_functional.PromptExecutionSettings
	ChatHistory           ai_functional.ChatHistory
	Function              functions.KernelFunction
	Kernel                abstractions.Kernel
	Result                functions.FunctionResult
	Terminate             bool
}

type IAutoFunctionInvocationFilter interface {
	OnAutoFunctionInvocation(context AutoFunctionInvocationContext, next func(AutoFunctionInvocationContext) error) error
}
