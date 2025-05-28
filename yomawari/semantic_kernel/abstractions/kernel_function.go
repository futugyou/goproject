package abstractions

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
)

type KernelFunction interface {
	functions.AIFunction
	GetPluginName() string
	GetMetadata() KernelFunctionMetadata
	WithKernel(kernel *Kernel, pluginName *string) KernelFunction
	InvokeFunction(ctx context.Context, kernel Kernel, arguments KernelArguments) (*FunctionResult, error)
	InvokeStreaming(ctx context.Context, kernel Kernel, arguments KernelArguments) (<-chan StreamingKernelContent, <-chan error)
}
