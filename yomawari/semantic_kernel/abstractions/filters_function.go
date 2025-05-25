package abstractions

import (
	"context"
)

type FunctionInvocationContext struct {
	Ctx         context.Context
	CancelFunc  context.CancelFunc
	IsStreaming bool
	Kernel      Kernel
	Function    KernelFunction
	Arguments   KernelArguments
	Result      FunctionResult
}

type IFunctionInvocationFilter interface {
	OnFunctionInvocation(context FunctionInvocationContext, next func(FunctionInvocationContext) error) error
}
