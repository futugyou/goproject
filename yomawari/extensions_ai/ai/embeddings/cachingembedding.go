package embeddings

import (
	"context"
	"fmt"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/embeddings"
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

func (g *CachingEmbeddingGenerator[TInput, TEmbedding]) Generate(ctx context.Context, values []TInput, options *embeddings.EmbeddingGenerationOptions) (*embeddings.GeneratedEmbeddings[TEmbedding], error) {
	if len(values) == 0 {
		return embeddings.NewGeneratedEmbeddings[TEmbedding](), nil
	}

	if len(values) == 1 {
		var cacheKey = g.GetCacheKey(values[0], options)
		if cachedResponse, err := g.ReadCache(ctx, cacheKey); err == nil {
			return embeddings.NewGeneratedEmbeddingsFromCollection[TEmbedding]([]TEmbedding{*cachedResponse}), nil
		}

		response, err := g.InnerGenerator.Generate(ctx, values, options)
		if err != nil {
			return nil, err
		}

		if response.Count() != 1 {
			return nil, fmt.Errorf("expected exactly one embedding to be generated, but received %d", response.Count())
		}

		g.WriteCache(ctx, cacheKey, response.Get(0))
		return response, err
	}

	var results = make([]TEmbedding, len(values))
	var uncached []struct {
		index    int
		cacheKey string
		input    TInput
	}

	for i, value := range values {
		var cacheKey = g.GetCacheKey(value, options)
		if existing, err := g.ReadCache(ctx, cacheKey); err == nil {
			results[i] = *existing
		} else {
			uncached = append(uncached, struct {
				index    int
				cacheKey string
				input    TInput
			}{i, cacheKey, value})
		}
	}

	if len(uncached) > 0 {
		inputs := make([]TInput, len(uncached))
		for i, u := range uncached {
			inputs[i] = u.input
		}

		uncachedResults, err := g.Generate(ctx, inputs, options)
		if err != nil {
			return nil, err
		}

		for i, u := range uncached {
			g.WriteCache(ctx, u.cacheKey, uncachedResults.Get(i))
			results[u.index] = uncachedResults.Get(i)
		}
	}

	return embeddings.NewGeneratedEmbeddingsFromCollection[TEmbedding](results), nil
}
