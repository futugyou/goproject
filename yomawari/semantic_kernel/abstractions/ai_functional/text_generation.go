package ai_functional

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/services"
)

type ITextGenerationService interface {
	services.IAIService
	GetTextContents(ctx context.Context, prompt string, executionSettings *PromptExecutionSettings, kernel *abstractions.Kernel) ([]contents.TextContent, error)
	GetTextContent(ctx context.Context, prompt string, executionSettings *PromptExecutionSettings, kernel *abstractions.Kernel) (*contents.TextContent, error)
	GetStreamingTextContents(ctx context.Context, prompt string, executionSettings *PromptExecutionSettings, kernel *abstractions.Kernel) (<-chan contents.StreamingTextContent, <-chan error)
}
