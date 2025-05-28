package abstractions

import (
	"context"
	"sync"

	"github.com/futugyou/yomawari/core"
)

type Kernel struct {
	data                          map[string]any
	plugins                       KernelPluginCollection
	functionInvocationFilters     []IFunctionInvocationFilter
	promptRenderFilters           []IPromptRenderFilter
	autoFunctionInvocationFilters []IAutoFunctionInvocationFilter
	Services                      core.IServiceProvider
}

var once sync.Once

func NewKernel(services core.IServiceProvider, plugins *KernelPluginCollection) *Kernel {
	if services == nil {
		services = &core.ServiceProvider{}
	}
	if plugins == nil {
		plugins, _ = core.GetService[*KernelPluginCollection](services)
	}
	if plugins == nil {
		plugins = &KernelPluginCollection{}
		pluginss, _ := core.GetService[[]KernelPlugin](services)
		if len(pluginss) > 0 {
			plugins.AddRange(pluginss)
		}
	}
	kernel := &Kernel{
		data:     make(map[string]any),
		Services: services,
	}

	if plugins != nil {
		kernel.plugins = *plugins
	}

	functionInvocationFilters, _ := core.GetService[[]IFunctionInvocationFilter](services)
	if len(functionInvocationFilters) > 0 {
		kernel.functionInvocationFilters = functionInvocationFilters
	}
	promptRenderFilters, _ := core.GetService[[]IPromptRenderFilter](services)
	if len(promptRenderFilters) > 0 {
		kernel.promptRenderFilters = promptRenderFilters
	}
	autoFunctionInvocationFilters, _ := core.GetService[[]IAutoFunctionInvocationFilter](services)
	if len(autoFunctionInvocationFilters) > 0 {
		kernel.autoFunctionInvocationFilters = autoFunctionInvocationFilters
	}
	return kernel
}

func (k *Kernel) GetFunctionInvocationFilters() []IFunctionInvocationFilter {
	once.Do(func() {
		if k.functionInvocationFilters == nil {
			k.functionInvocationFilters = []IFunctionInvocationFilter{}
		}
	})
	return k.functionInvocationFilters
}

func (k *Kernel) GetPromptRenderFilters() []IPromptRenderFilter {
	once.Do(func() {
		if k.promptRenderFilters == nil {
			k.promptRenderFilters = []IPromptRenderFilter{}
		}
	})
	return k.promptRenderFilters
}

func (k *Kernel) GetAutoFunctionInvocationFilters() []IAutoFunctionInvocationFilter {
	once.Do(func() {
		if k.autoFunctionInvocationFilters == nil {
			k.autoFunctionInvocationFilters = []IAutoFunctionInvocationFilter{}
		}
	})
	return k.autoFunctionInvocationFilters
}

func (k *Kernel) GetPlugins() KernelPluginCollection {
	once.Do(func() {
		if k.plugins.Dictionary == nil {
			k.plugins.Dictionary = core.NewDictionary[string, KernelPlugin]()
		}
	})
	return k.plugins
}

func (k *Kernel) GetDatas() map[string]any {
	once.Do(func() {
		if k.data == nil {
			k.data = map[string]any{}
		}
	})
	return k.data
}

func (k *Kernel) OnFunctionInvocation(ctx context.Context, function KernelFunction, arguments KernelArguments, functionResult FunctionResult, isStreaming bool, functionCallback func(FunctionInvocationContext) error) (*FunctionInvocationContext, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	con := &FunctionInvocationContext{
		Ctx:         ctx,
		CancelFunc:  cancelFunc,
		IsStreaming: isStreaming,
		Kernel:      *k,
		Function:    function,
		Arguments:   arguments,
		Result:      functionResult,
	}
	return con, k.invokeFilterOrFunction(k.functionInvocationFilters, functionCallback, *con, 0)
}

func (k *Kernel) invokeFilterOrFunction(functionFilters []IFunctionInvocationFilter, functionCallback func(FunctionInvocationContext) error, context FunctionInvocationContext, index int) error {
	if len(functionFilters) > 0 && index < len(functionFilters) {
		next := func(context FunctionInvocationContext) error {
			return k.invokeFilterOrFunction(functionFilters, functionCallback, context, index+1)
		}
		return functionFilters[index].OnFunctionInvocation(context, next)
	} else {
		return functionCallback(context)
	}
}

func (k *Kernel) OnPromptRender(ctx context.Context, function KernelFunction, arguments KernelArguments, isStreaming bool, executionSettings PromptExecutionSettings, functionCallback func(PromptRenderContext) error) (*PromptRenderContext, error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	con := &PromptRenderContext{
		Ctx:               ctx,
		CancelFunc:        cancelFunc,
		IsStreaming:       isStreaming,
		Kernel:            *k,
		Function:          function,
		Arguments:         arguments,
		ExecutionSettings: executionSettings,
	}
	return con, k.invokeFilterOrPromptRender(k.promptRenderFilters, functionCallback, *con, 0)
}

func (k *Kernel) invokeFilterOrPromptRender(functionFilters []IPromptRenderFilter, functionCallback func(PromptRenderContext) error, context PromptRenderContext, index int) error {
	if len(functionFilters) > 0 && index < len(functionFilters) {
		next := func(context PromptRenderContext) error {
			return k.invokeFilterOrPromptRender(functionFilters, functionCallback, context, index+1)
		}
		return functionFilters[index].OnPromptRender(context, next)
	} else {
		return functionCallback(context)
	}
}

func (k *Kernel) Invoke(ctx context.Context, pluginName string, functionName string, arguments KernelArguments) (*FunctionResult, error) {
	function, err := k.plugins.GetFunction(pluginName, functionName)
	if err != nil {
		return nil, err
	}

	return function.InvokeFunction(ctx, *k, arguments)
}

func (k *Kernel) InvokeStreaming(ctx context.Context, pluginName string, functionName string, arguments KernelArguments) (<-chan StreamingKernelContent, <-chan error) {
	resultCh := make(chan StreamingKernelContent)
	errCh := make(chan error, 1)
	defer close(resultCh)
	defer close(errCh)
	function, err := k.plugins.GetFunction(pluginName, functionName)
	if err != nil {
		errCh <- err
		return resultCh, errCh
	}

	return function.InvokeStreaming(ctx, *k, arguments)
}
