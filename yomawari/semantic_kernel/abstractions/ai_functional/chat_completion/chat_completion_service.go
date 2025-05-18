package chat_completion

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/ai_functional"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/services"
)

type IChatCompletionService interface {
	services.IAIService
	GetChatMessageContents(ctx context.Context, chatHistory ChatHistory, executionSettings ai_functional.PromptExecutionSettings, kernel abstractions.Kernel) ([]contents.ChatMessageContent, error)
	GetStreamingChatMessageContents(ctx context.Context, chatHistory ChatHistory, executionSettings ai_functional.PromptExecutionSettings, kernel abstractions.Kernel) (<-chan contents.StreamingChatMessageContent, <-chan error)
}
