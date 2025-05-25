package abstractions

import (
	"context"
)

type IAudioToText interface {
	IAIService
	GetTextContents(ctx context.Context, content AudioContent, executionSettings *PromptExecutionSettings, kernel *Kernel) ([]TextContent, error)
	GetTextContent(ctx context.Context, content AudioContent, executionSettings *PromptExecutionSettings, kernel *Kernel) (*TextContent, error)
}
