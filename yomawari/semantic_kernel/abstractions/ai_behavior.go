package abstractions

import (
	"fmt"

	"github.com/futugyou/yomawari/semantic_kernel/utilities"
)

const FunctionNameSeparator = "."

type FunctionChoiceBehavior struct {
	DefaultOptions FunctionChoiceBehaviorOptions `json:"options"`
	Functions      []string                      `json:"functions"`
	functions      []KernelFunction              `json:"-"`
	behaviorType   string                        `json:"-"`
	autoInvoke     bool
}

func NewFunctionChoiceBehavior(functions []KernelFunction, options *FunctionChoiceBehaviorOptions, behaviorType string, autoInvoke bool) *FunctionChoiceBehavior {
	if options == nil {
		options = &FunctionChoiceBehaviorOptions{}
	}
	fs := []string{}
	for _, v := range functions {
		fs = append(fs, utilities.ToFullyQualifiedName(v.GetName(), v.GetPluginName(), FunctionNameSeparator))
	}
	return &FunctionChoiceBehavior{
		DefaultOptions: *options,
		functions:      functions,
		Functions:      fs,
		behaviorType:   behaviorType,
		autoInvoke:     autoInvoke,
	}
}

var trueBehavior = true

func AutoFunctionChoiceBehavior(functions []KernelFunction, options *FunctionChoiceBehaviorOptions, autoInvoke *bool) *FunctionChoiceBehavior {
	if autoInvoke == nil {
		autoInvoke = &trueBehavior
	}
	return NewFunctionChoiceBehavior(functions, options, "auto", *autoInvoke)
}

func RequiredFunctionChoiceBehavior(functions []KernelFunction, options *FunctionChoiceBehaviorOptions, autoInvoke *bool) *FunctionChoiceBehavior {
	if autoInvoke == nil {
		autoInvoke = &trueBehavior
	}
	return NewFunctionChoiceBehavior(functions, options, "required", *autoInvoke)
}

func NoneFunctionChoiceBehavior(functions []KernelFunction, options *FunctionChoiceBehaviorOptions) *FunctionChoiceBehavior {
	return NewFunctionChoiceBehavior(functions, options, "none", false)
}

func (f *FunctionChoiceBehavior) GetBehaviorType() string {
	if f == nil {
		return ""
	}
	return f.behaviorType
}

func (f *FunctionChoiceBehavior) GetConfiguration(context FunctionChoiceBehaviorConfigurationContext) FunctionChoiceBehaviorConfiguration {
	if f.behaviorType == "required" {
		if context.RequestSequenceIndex >= 1 {
			return FunctionChoiceBehaviorConfiguration{
				Options:    f.DefaultOptions,
				Choice:     FunctionChoiceRequired,
				Functions:  nil,
				AutoInvoke: f.autoInvoke,
			}
		}
	}

	choice := FunctionChoiceAuto
	if f.behaviorType == "none" {
		f.autoInvoke = false
		choice = FunctionChoiceNone
	}

	fs, err := f.GetFunctions(f.Functions, context.Kernel, f.autoInvoke)
	if err != nil {
		fs = []KernelFunction{}
	}

	return FunctionChoiceBehaviorConfiguration{
		Options:    f.DefaultOptions,
		Choice:     choice,
		Functions:  fs,
		AutoInvoke: f.autoInvoke,
	}
}

func (f *FunctionChoiceBehavior) GetFunctions(functionFQNs []string, kernel *Kernel, autoInvoke bool) ([]KernelFunction, error) {
	if autoInvoke && kernel == nil {
		return nil, fmt.Errorf("auto-invocation is not supported when no kernel is provided")
	}

	availableFunctions := []KernelFunction{}

	if len(functionFQNs) > 0 {
		for _, functionFQN := range functionFQNs {
			nameParts := utilities.ParseFunctionName(functionFQN, FunctionNameSeparator)
			if kernel != nil {
				plugins := kernel.GetPlugins()
				if function, err := plugins.GetFunction(nameParts.PluginName, nameParts.Name); err != nil {
					availableFunctions = append(availableFunctions, function)
					continue
				}
			}

			// If auto-invocation is requested and no function is found in the kernel, fail early.
			if autoInvoke {
				return nil, fmt.Errorf("the specified function %s is not available in the kernel", functionFQN)
			}

			var function *KernelFunction
			for _, v := range f.functions {
				if v.GetName() == nameParts.Name && v.GetPluginName() == nameParts.PluginName {
					function = &v
					break
				}
			}
			if function != nil {
				availableFunctions = append(availableFunctions, *function)
				continue
			}

			return nil, fmt.Errorf("the specified function %s was not found", functionFQN)
		}
	} else if len(functionFQNs) == 0 {
		return availableFunctions, nil
	} else if kernel != nil {
		for _, plugin := range kernel.GetPlugins().Items() {
			availableFunctions = append(availableFunctions, plugin.GetFunctions()...)
		}
	}

	return availableFunctions, nil
}
