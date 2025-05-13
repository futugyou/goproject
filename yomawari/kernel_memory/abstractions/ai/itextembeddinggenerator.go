package ai

import "context"

type ITextEmbeddingBatchGenerator interface {
	GetMaxBatchSize() int64
	GenerateEmbeddingBatch(ctx context.Context, textList []string) ([]Embedding, error)
}

type ITextEmbeddingGenerator interface {
	ITextTokenizer
	GetMaxTokens() int64
	GenerateEmbedding(ctx context.Context, text string) (Embedding, error)
}
