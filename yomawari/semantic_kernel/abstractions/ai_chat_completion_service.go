package abstractions

import (
	"context"
)

type IChatCompletionService interface {
	IAIService
	GetChatMessageContents(ctx context.Context, chatHistory ChatHistory, executionSettings *PromptExecutionSettings, kernel *Kernel) ([]ChatMessageContent, error)
	GetStreamingChatMessageContents(ctx context.Context, chatHistory ChatHistory, executionSettings *PromptExecutionSettings, kernel *Kernel) (<-chan StreamingChatMessageContent, <-chan error)
}
