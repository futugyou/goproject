package chatcompletion

import (
	"context"
	"encoding/json"

	"github.com/futugyou/ai-extension/abstractions/chatcompletion"
	"github.com/futugyou/ai-extension/core"
)

type DistributedCachingChatClient struct {
	*CachingChatClient
	storage core.IDistributedCache
}

func NewDistributedCachingChatClient(
	innerClient chatcompletion.IChatClient,
	storage core.IDistributedCache,
) *DistributedCachingChatClient {
	return &DistributedCachingChatClient{
		CachingChatClient: NewCachingChatClient(innerClient),
		storage:           storage,
	}
}

func (client *DistributedCachingChatClient) ReadCache(ctx context.Context, key string) (*chatcompletion.ChatResponse, error) {
	cache, err := client.storage.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var result chatcompletion.ChatResponse
	if err := json.Unmarshal(cache, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *DistributedCachingChatClient) ReadCacheStreaming(ctx context.Context, key string) ([]chatcompletion.ChatResponseUpdate, error) {
	cache, err := client.storage.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var result []chatcompletion.ChatResponseUpdate
	if err := json.Unmarshal(cache, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (client *DistributedCachingChatClient) WriteCache(ctx context.Context, key string, value chatcompletion.ChatResponse) error {
	cache, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.storage.Set(ctx, key, cache)
}

func (client *DistributedCachingChatClient) WriteCacheStreaming(ctx context.Context, key string, value []chatcompletion.ChatResponseUpdate) error {
	cache, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.storage.Set(ctx, key, cache)
}

func (client *DistributedCachingChatClient) GetCacheKey(boxed bool, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) string {
	panic("GetCacheKey must be implemented by subclass")
}
