package abstractions

import (
	"fmt"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
)

type KernelPlugin interface {
	Name() string
	Description() string
	FunctionCount() int
	GetFunction(name string) (KernelFunction, error)
	GetFunctions() []KernelFunction
	TryGetFunction(name string) (KernelFunction, bool)
	Contains(name string) bool
	ContainsFunc(function KernelFunction) bool
	GetFunctionsMetadata() []KernelFunctionMetadata
	AsAIFunctions(kernel *Kernel) []functions.AIFunction
	IterateFunctions() []KernelFunction
}

type BaseKernelPlugin struct {
	name        string
	description string
	functions   map[string]KernelFunction
}

func NewBaseKernelPlugin(name, description string) *BaseKernelPlugin {
	if name == "" {
		panic("plugin name cannot be empty")
	}
	return &BaseKernelPlugin{
		name:        name,
		description: description,
		functions:   make(map[string]KernelFunction),
	}
}

func (p *BaseKernelPlugin) Name() string {
	return p.name
}

func (p *BaseKernelPlugin) Description() string {
	return p.description
}

func (p *BaseKernelPlugin) FunctionCount() int {
	return len(p.functions)
}

func (p *BaseKernelPlugin) GetFunction(name string) (KernelFunction, error) {
	if fn, ok := p.functions[name]; ok {
		return fn, nil
	}
	return nil, fmt.Errorf("function '%s' not found in plugin '%s'", name, p.name)
}

func (p *BaseKernelPlugin) GetFunctions() []KernelFunction {
	result := []KernelFunction{}
	for _, v := range p.functions {
		result = append(result, v)
	}
	return result
}

func (p *BaseKernelPlugin) TryGetFunction(name string) (KernelFunction, bool) {
	fn, ok := p.functions[name]
	return fn, ok
}

func (p *BaseKernelPlugin) Contains(name string) bool {
	_, ok := p.functions[name]
	return ok
}

func (p *BaseKernelPlugin) ContainsFunc(function KernelFunction) bool {
	return p.Contains(function.GetName())
}

func (p *BaseKernelPlugin) GetFunctionsMetadata() []KernelFunctionMetadata {
	meta := make([]KernelFunctionMetadata, 0, len(p.functions))
	for _, fn := range p.functions {
		meta = append(meta, fn.GetMetadata())
	}
	return meta
}

func (p *BaseKernelPlugin) AsAIFunctions(kernel *Kernel) []functions.AIFunction {
	var result []functions.AIFunction
	for _, fn := range p.functions {
		result = append(result, fn.WithKernel(kernel, nil))
	}
	return result
}

func (p *BaseKernelPlugin) IterateFunctions() []KernelFunction {
	result := make([]KernelFunction, 0, len(p.functions))
	for _, fn := range p.functions {
		result = append(result, fn)
	}
	return result
}

func RegisterPlugin(name string, description string, funcs ...KernelFunction) KernelPlugin {
	base := NewBaseKernelPlugin(name, description)
	for _, f := range funcs {
		base.functions[f.GetName()] = f
	}
	return base
}
