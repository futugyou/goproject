package ai_functional

import (
	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
)

type FunctionChoiceBehaviorConfigurationContext struct {
	ChatHistory          ChatHistory
	Kernel               abstractions.Kernel
	RequestSequenceIndex int
}
