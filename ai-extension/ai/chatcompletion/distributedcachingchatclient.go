package chatcompletion

import (
	"context"
	"encoding/json"

	"github.com/futugyou/ai-extension/abstractions/chatcompletion"
	"github.com/futugyou/ai-extension/core"
)

type DistributedCachingChatClient struct {
	chatcompletion.DelegatingChatClient
	storage core.IDistributedCache
}

func NewDistributedCachingChatClient(
	innerClient chatcompletion.IChatClient,
	storage core.IDistributedCache,
) *DistributedCachingChatClient {
	return &DistributedCachingChatClient{
		DelegatingChatClient: chatcompletion.DelegatingChatClient{
			InnerClient: innerClient,
		},
		storage: storage,
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
