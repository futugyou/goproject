package core

import "github.com/futugyou/yomawari/semantic_kernel/abstractions"

type TextSearchProviderOptions struct {
	Top                       int
	SearchTime                RagBehavior
	PluginFunctionName        string
	PluginFunctionDescription string
	ContextPrompt             string
	IncludeCitationsPrompt    string
	ContextFormatter          func(context []abstractions.TextSearchResult) string
}

type RagBehavior string

const (
	RagBehaviorBeforeAIInvoke          = "BeforeAIInvoke"
	RagBehaviorOnDemandFunctionCalling = "OnDemandFunctionCalling"
)
