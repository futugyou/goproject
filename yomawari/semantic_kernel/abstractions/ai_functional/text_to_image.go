package ai_functional

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/services"
)

type ITextToImageService interface {
	services.IAIService
	GetImageContents(ctx context.Context, input contents.TextContent, executionSettings *PromptExecutionSettings, kernel *abstractions.Kernel) ([]contents.ImageContent, error)
	GetImageContent(ctx context.Context, input contents.TextContent, executionSettings *PromptExecutionSettings, kernel *abstractions.Kernel) (*contents.ImageContent, error)
}
