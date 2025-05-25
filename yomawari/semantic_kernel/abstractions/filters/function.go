package filters

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/functions"
)

type FunctionInvocationContext struct {
	Ctx         context.Context
	CancelFunc  context.CancelFunc
	IsStreaming bool
	Kernel      abstractions.Kernel
	Function    functions.KernelFunction
	Arguments   functions.KernelArguments
	Result      functions.FunctionResult
}

type IFunctionInvocationFilter interface {
	OnFunctionInvocation(context FunctionInvocationContext, next func(FunctionInvocationContext) error) error
}
