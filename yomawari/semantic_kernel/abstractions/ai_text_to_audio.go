package abstractions

import (
	"context"
)

type ITextToAudioService interface {
	IAIService
	GetAudioContents(ctx context.Context, prompt string, executionSettings *PromptExecutionSettings, kernel *Kernel) ([]AudioContent, error)
	GetAudioContent(ctx context.Context, prompt string, executionSettings *PromptExecutionSettings, kernel *Kernel) (*AudioContent, error)
}
