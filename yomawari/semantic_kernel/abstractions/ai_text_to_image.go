package abstractions

import (
	"context"
)

type ITextToImageService interface {
	IAIService
	GetImageContents(ctx context.Context, input TextContent, executionSettings *PromptExecutionSettings, kernel *Kernel) ([]ImageContent, error)
	GetImageContent(ctx context.Context, input TextContent, executionSettings *PromptExecutionSettings, kernel *Kernel) (*ImageContent, error)
}
