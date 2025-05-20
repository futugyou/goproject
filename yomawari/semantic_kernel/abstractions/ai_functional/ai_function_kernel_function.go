package ai_functional

import (
	aifunctions "github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/functions"
)

func AIFunctionTosKernelFunction(fn aifunctions.AIFunction) functions.KernelFunction {
	//TODO: need KernelFunction first
	// return NewAIFunctionKernelFunctionFromAIFunction(fn)
	panic("unimplemented")
}

type AIFunctionKernelFunction struct {
}
