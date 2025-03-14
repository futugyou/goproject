package embeddings

import (
	"context"
	"encoding/json"

	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/generative-ai/abstractions/embeddings"
	"github.com/futugyou/yomawari/generative-ai/abstractions/utilities"
)

type DistributedCachingEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding] struct {
	*CachingEmbeddingGenerator[TInput, TEmbedding]
	storage core.IDistributedCache
}

func NewDistributedCachingEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding](
	innerGenerator embeddings.IEmbeddingGenerator[TInput, TEmbedding],
	storage core.IDistributedCache,
) *DistributedCachingEmbeddingGenerator[TInput, TEmbedding] {
	return &DistributedCachingEmbeddingGenerator[TInput, TEmbedding]{
		CachingEmbeddingGenerator: NewCachingEmbeddingGenerator[TInput, TEmbedding](innerGenerator),
		storage:                   storage,
	}
}

func (client *DistributedCachingEmbeddingGenerator[TInput, TEmbedding]) ReadCache(ctx context.Context, key string) (*TEmbedding, error) {
	cache, err := client.storage.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var result TEmbedding
	if err := json.Unmarshal(cache, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *DistributedCachingEmbeddingGenerator[TInput, TEmbedding]) WriteCache(ctx context.Context, key string, value TEmbedding) error {
	cache, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.storage.Set(ctx, key, cache)
}

func (client *DistributedCachingEmbeddingGenerator[TInput, TEmbedding]) GetCacheKey(values ...interface{}) string {
	return utilities.HashDataToString(values...)
}
