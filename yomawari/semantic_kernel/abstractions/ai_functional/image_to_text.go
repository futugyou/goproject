package ai_functional

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/services"
)

type IImageToTextService interface {
	services.IAIService
	GetTextContents(ctx context.Context, content contents.ImageContent, executionSettings *PromptExecutionSettings, kernel *abstractions.Kernel) ([]contents.TextContent, error)
	GetTextContent(ctx context.Context, content contents.ImageContent, executionSettings *PromptExecutionSettings, kernel *abstractions.Kernel) (*contents.TextContent, error)
}
