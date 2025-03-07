package embeddings

import (
	"context"

	"github.com/futugyou/ai-extension/abstractions/embeddings"
)

type CachingEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding] struct {
	embeddings.DelegatingEmbeddingGenerator[TInput, TEmbedding]
}

func NewCachingEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding](
	innerGenerator embeddings.IEmbeddingGenerator[TInput, TEmbedding],
) *CachingEmbeddingGenerator[TInput, TEmbedding] {
	return &CachingEmbeddingGenerator[TInput, TEmbedding]{
		DelegatingEmbeddingGenerator: *embeddings.NewDelegatingEmbeddingGenerator[TInput, TEmbedding](innerGenerator),
	}
}

func (client *CachingEmbeddingGenerator[TInput, TEmbedding]) GetCacheKey(values ...interface{}) string {
	panic("GetCacheKey must be implemented by subclass")
}

func (client *CachingEmbeddingGenerator[TInput, TEmbedding]) ReadCache(ctx context.Context, key string) (*TEmbedding, error) {
	panic("ReadCache must be implemented by subclass")
}

func (client *CachingEmbeddingGenerator[TInput, TEmbedding]) WriteCache(ctx context.Context, key string, value TEmbedding) error {
	panic("WriteCache must be implemented by subclass")
}
