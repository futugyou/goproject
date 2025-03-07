package embeddings

import "context"

type DelegatingEmbeddingGenerator[TInput any, TEmbedding IEmbedding] struct {
	InnerGenerator IEmbeddingGenerator[TInput, TEmbedding]
}

func NewDelegatingEmbeddingGenerator[TInput any, TEmbedding IEmbedding](innerGenerator IEmbeddingGenerator[TInput, TEmbedding]) *DelegatingEmbeddingGenerator[TInput, TEmbedding] {
	return &DelegatingEmbeddingGenerator[TInput, TEmbedding]{
		InnerGenerator: innerGenerator,
	}
}

func (g *DelegatingEmbeddingGenerator[TInput, TEmbedding]) Generate(ctx context.Context, values []TInput, options *EmbeddingGenerationOptions) (*GeneratedEmbeddings[TEmbedding], error) {
	return g.InnerGenerator.Generate(ctx, values, options)
}
