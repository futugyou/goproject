package chatcompletion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/utilities"
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
	if client.storage == nil {
		return nil, fmt.Errorf("storage is nil")
	}
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
	if client.storage == nil {
		return nil, fmt.Errorf("storage is nil")
	}
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
	if client.storage == nil {
		return fmt.Errorf("storage is nil")
	}
	cache, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.storage.Set(ctx, key, cache)
}

func (client *DistributedCachingChatClient) WriteCacheStreaming(ctx context.Context, key string, value []chatcompletion.ChatResponseUpdate) error {
	if client.storage == nil {
		return fmt.Errorf("storage is nil")
	}
	cache, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.storage.Set(ctx, key, cache)
}

func (client *DistributedCachingChatClient) GetCacheKey(values ...interface{}) string {
	return utilities.HashDataToString(values...)
}
