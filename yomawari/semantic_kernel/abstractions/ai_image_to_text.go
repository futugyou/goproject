package abstractions

import (
	"context"
)

type IImageToTextService interface {
	IAIService
	GetTextContents(ctx context.Context, content ImageContent, executionSettings *PromptExecutionSettings, kernel *Kernel) ([]TextContent, error)
	GetTextContent(ctx context.Context, content ImageContent, executionSettings *PromptExecutionSettings, kernel *Kernel) (*TextContent, error)
}
