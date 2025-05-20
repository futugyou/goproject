package ai_functional

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/services"
)

type IChatCompletionService interface {
	services.IAIService
	GetChatMessageContents(ctx context.Context, chatHistory ChatHistory, executionSettings PromptExecutionSettings, kernel abstractions.Kernel) ([]contents.ChatMessageContent, error)
	GetStreamingChatMessageContents(ctx context.Context, chatHistory ChatHistory, executionSettings PromptExecutionSettings, kernel abstractions.Kernel) (<-chan contents.StreamingChatMessageContent, <-chan error)
}
