package abstractions

import (
	"context"
)

type ITextGenerationService interface {
	IAIService
	GetTextContents(ctx context.Context, prompt string, executionSettings *PromptExecutionSettings, kernel *Kernel) ([]TextContent, error)
	GetTextContent(ctx context.Context, prompt string, executionSettings *PromptExecutionSettings, kernel *Kernel) (*TextContent, error)
	GetStreamingTextContents(ctx context.Context, prompt string, executionSettings *PromptExecutionSettings, kernel *Kernel) (<-chan StreamingTextContent, <-chan error)
}
