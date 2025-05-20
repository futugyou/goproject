package ai_functional

import (
	"fmt"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/functions"
	"github.com/futugyou/yomawari/semantic_kernel/utilities"
)

const FunctionNameSeparator = "."

type FunctionChoiceBehavior struct {
	DefaultOptions FunctionChoiceBehaviorOptions `json:"options"`
	Functions      []string                      `json:"functions"`
	functions      []functions.KernelFunction    `json:"-"`
	behaviorType   string                        `json:"-"`
	autoInvoke     bool
}

func NewFunctionChoiceBehavior(functions []functions.KernelFunction, options *FunctionChoiceBehaviorOptions, behaviorType string, autoInvoke bool) *FunctionChoiceBehavior {
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

	fs, err := f.GetFunctions(f.Functions, &context.Kernel, f.autoInvoke)
	if err != nil {
		fs = []functions.KernelFunction{}
	}

	return FunctionChoiceBehaviorConfiguration{
		Options:    f.DefaultOptions,
		Choice:     choice,
		Functions:  fs,
		AutoInvoke: f.autoInvoke,
	}
}

func (f *FunctionChoiceBehavior) GetFunctions(functionFQNs []string, kernel *abstractions.Kernel, autoInvoke bool) ([]functions.KernelFunction, error) {
	if autoInvoke && kernel == nil {
		return nil, fmt.Errorf("auto-invocation is not supported when no kernel is provided")
	}

	availableFunctions := []functions.KernelFunction{}

	if len(functionFQNs) > 0 {
		for _, functionFQN := range functionFQNs {
			_ = utilities.ParseFunctionName(functionFQN, FunctionNameSeparator)
			// TODO: A complete kernel definition is required
			if kernel != nil {
				// if function, ok := kernel.Plugins.TryGetFunction(nameParts.PluginName, nameParts.Name); ok {
				// 	availableFunctions = append(availableFunctions, function)
				// 	continue
				// }
			}

			// If auto-invocation is requested and no function is found in the kernel, fail early.
			if autoInvoke {
				return nil, fmt.Errorf("the specified function %s is not available in the kernel", functionFQN)
			}

			// // TODO: A complete KernelFunction definition is required
			var function *functions.KernelFunction
			for _, _ = range f.Functions {
				// if v.Name == nameParts.Name && v.PluginName == nameParts.PluginName {
				// 	function = &v
				// 	break
				// }
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
		// for _, plugin := range kernel.Plugins {
		// 	availableFunctions = append(availableFunctions, plugin)
		// }
	}

	return availableFunctions, nil
}
