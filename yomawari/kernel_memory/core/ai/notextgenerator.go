package ai

import (
	"context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
)

type NoTextGenerator struct {
}

// CountTokens implements ai.ITextGenerator.
func (n *NoTextGenerator) CountTokens(ctx context.Context, text string) int64 {
	panic("text generation has been disabled")
}

// GenerateText implements ai.ITextGenerator.
func (n *NoTextGenerator) GenerateText(ctx context.Context, prompt string, options *ai.TextGenerationOptions) <-chan ai.GenerateTextResponse {
	panic("text generation has been disabled")
}

// GetMaxTokenTotal implements ai.ITextGenerator.
func (n *NoTextGenerator) GetMaxTokenTotal() int64 {
	panic("text generation has been disabled")
}

// GetTokens implements ai.ITextGenerator.
func (n *NoTextGenerator) GetTokens(ctx context.Context, text string) []string {
	panic("text generation has been disabled")
}

var _ ai.ITextGenerator = (*NoTextGenerator)(nil)
