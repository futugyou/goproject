package ai

import (
	"context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
)

type NoEmbeddingGenerator struct {
}

// CountTokens implements ai.ITextEmbeddingGenerator.
func (n *NoEmbeddingGenerator) CountTokens(ctx context.Context, text string) int64 {
	panic("embedding generation has been disabled")
}

// GetTokens implements ai.ITextEmbeddingGenerator.
func (n *NoEmbeddingGenerator) GetTokens(ctx context.Context, text string) []string {
	panic("embedding generation has been disabled")
}

// GenerateEmbedding implements ai.ITextEmbeddingGenerator.
func (n *NoEmbeddingGenerator) GenerateEmbedding(ctx context.Context, text string) (ai.Embedding, error) {
	panic("embedding generation has been disabled")
}

// GetMaxTokens implements ai.ITextEmbeddingGenerator.
func (n *NoEmbeddingGenerator) GetMaxTokens() int64 {
	panic("embedding generation has been disabled")
}

var _ ai.ITextEmbeddingGenerator = (*NoEmbeddingGenerator)(nil)
