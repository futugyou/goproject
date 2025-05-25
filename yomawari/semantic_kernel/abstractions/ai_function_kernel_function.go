package abstractions

import (
	aifunctions "github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
)

func AIFunctionTosKernelFunction(fn aifunctions.AIFunction) KernelFunction {
	//TODO: need KernelFunction first
	// return NewAIFunctionKernelFunctionFromAIFunction(fn)
	panic("unimplemented")
}

type AIFunctionKernelFunction struct {
	pluginName string
	aiFunction aifunctions.AIFunction
}

// GetAdditionalProperties implements functions.KernelFunction.
func (a *AIFunctionKernelFunction) GetAdditionalProperties() map[string]interface{} {
	panic("unimplemented")
}

// GetDescription implements functions.KernelFunction.
func (a *AIFunctionKernelFunction) GetDescription() string {
	return a.aiFunction.GetName()
}

// GetName implements functions.KernelFunction.
func (a *AIFunctionKernelFunction) GetName() string {
	panic("unimplemented")
}

// GetPluginName implements functions.KernelFunction.
func (a *AIFunctionKernelFunction) GetPluginName() string {
	panic("unimplemented")
}
