package ai_functional

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/services"
)

type ITextToAudioService interface {
	services.IAIService
	GetAudioContents(ctx context.Context, prompt string, executionSettings *PromptExecutionSettings, kernel *abstractions.Kernel) ([]contents.AudioContent, error)
	GetAudioContent(ctx context.Context, prompt string, executionSettings *PromptExecutionSettings, kernel *abstractions.Kernel) (*contents.AudioContent, error)
}
