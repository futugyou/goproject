package abstractions

import (
	"context"
)

type AutoFunctionInvocationContext struct {
	Ctx                   context.Context
	CancelFunc            context.CancelFunc
	IsStreaming           bool
	Arguments             KernelArguments
	RequestSequenceIndex  int
	FunctionSequenceIndex int
	FunctionCount         int
	ToolCallId            string
	ChatMessageContent    ChatMessageContent
	ExecutionSettings     PromptExecutionSettings
	ChatHistory           ChatHistory
	Function              KernelFunction
	Kernel                Kernel
	Result                FunctionResult
	Terminate             bool
}

type IAutoFunctionInvocationFilter interface {
	OnAutoFunctionInvocation(context AutoFunctionInvocationContext, next func(AutoFunctionInvocationContext) error) error
}
