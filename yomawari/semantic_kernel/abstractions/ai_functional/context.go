package ai_functional

import (
	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/functions"
)

type FunctionChoiceBehaviorConfigurationContext struct {
	ChatHistory          ChatHistory
	Kernel               abstractions.Kernel
	RequestSequenceIndex int
}

type FunctionChoiceBehaviorConfiguration struct {
	Choice     FunctionChoice
	Functions  []functions.KernelFunction
	AutoInvoke bool
	Options    FunctionChoiceBehaviorOptions
}
