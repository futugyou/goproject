package embeddings

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/embeddings"
)

type ConfigureOptionsEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding] struct {
	embeddings.DelegatingEmbeddingGenerator[TInput, TEmbedding]
	configureOptions func(*embeddings.EmbeddingGenerationOptions)
}

func NewConfigureOptionsEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding](
	innerGenerator embeddings.IEmbeddingGenerator[TInput, TEmbedding],
	configureOptions func(*embeddings.EmbeddingGenerationOptions),
) *ConfigureOptionsEmbeddingGenerator[TInput, TEmbedding] {
	return &ConfigureOptionsEmbeddingGenerator[TInput, TEmbedding]{
		DelegatingEmbeddingGenerator: *embeddings.NewDelegatingEmbeddingGenerator[TInput, TEmbedding](innerGenerator),
		configureOptions:             configureOptions,
	}
}

func (g *ConfigureOptionsEmbeddingGenerator[TInput, TEmbedding]) Generate(ctx context.Context, values []TInput, options *embeddings.EmbeddingGenerationOptions) (*embeddings.GeneratedEmbeddings[TEmbedding], error) {
	if options == nil {
		options = &embeddings.EmbeddingGenerationOptions{}
	}

	g.configureOptions(options)
	return g.DelegatingEmbeddingGenerator.Generate(ctx, values, options)
}
