package functions

import (
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
)

type KernelFunction interface {
	functions.AIFunction
	GetPluginName() string
	GetMetadata() KernelFunctionMetadata
	WithKernel(kernel *abstractions.Kernel, pluginName *string) KernelFunction
}
