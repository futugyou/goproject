package abstractions

import (
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
)

type KernelFunction interface {
	functions.AIFunction
	GetPluginName() string
	GetMetadata() KernelFunctionMetadata
	WithKernel(kernel *Kernel, pluginName *string) KernelFunction
}
