package abstractions

import "github.com/futugyou/yomawari/core"

type Kernel struct {
	data                          map[string]any
	plugins                       KernelPluginCollection
	functionInvocationFilters     []IFunctionInvocationFilter
	promptRenderFilters           []IPromptRenderFilter
	autoFunctionInvocationFilters []IAutoFunctionInvocationFilter
	Services                      core.IServiceProvider
}

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
