package embeddings

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/embeddings"
)

type GenerateFunc[TInput any, TEmbedding embeddings.IEmbedding] func(context.Context, []TInput, *embeddings.EmbeddingGenerationOptions, embeddings.IEmbeddingGenerator[TInput, TEmbedding]) (*embeddings.GeneratedEmbeddings[TEmbedding], error)

type AnonymousDelegatingEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding] struct {
	embeddings.DelegatingEmbeddingGenerator[TInput, TEmbedding]
	generateFunc GenerateFunc[TInput, TEmbedding]
}

func NewAnonymousDelegatingEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding](
	innerGenerator embeddings.IEmbeddingGenerator[TInput, TEmbedding],
	generateFunc GenerateFunc[TInput, TEmbedding],
) *AnonymousDelegatingEmbeddingGenerator[TInput, TEmbedding] {
	return &AnonymousDelegatingEmbeddingGenerator[TInput, TEmbedding]{
		DelegatingEmbeddingGenerator: *embeddings.NewDelegatingEmbeddingGenerator[TInput, TEmbedding](innerGenerator),
		generateFunc:                 generateFunc,
	}
}

func (g *AnonymousDelegatingEmbeddingGenerator[TInput, TEmbedding]) Generate(ctx context.Context, values []TInput, options *embeddings.EmbeddingGenerationOptions) (*embeddings.GeneratedEmbeddings[TEmbedding], error) {
	return g.generateFunc(ctx, values, options, &g.DelegatingEmbeddingGenerator)
}
