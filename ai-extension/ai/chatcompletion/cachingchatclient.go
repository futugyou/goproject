package chatcompletion

import (
	"context"

	"github.com/futugyou/ai-extension/abstractions/chatcompletion"
)

type CachingChatClient struct {
	chatcompletion.DelegatingChatClient
	CoalesceStreamingUpdates bool
}

func NewCachingChatClient(
	innerClient chatcompletion.IChatClient,
) *CachingChatClient {
	return &CachingChatClient{
		DelegatingChatClient: chatcompletion.DelegatingChatClient{
			InnerClient: innerClient,
		},
		CoalesceStreamingUpdates: true,
	}
}

func (client *CachingChatClient) GetResponse(
	ctx context.Context,
	chatMessages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
) (*chatcompletion.ChatResponse, error) {
	return client.InnerClient.GetResponse(ctx, chatMessages, options)
}

func (client *CachingChatClient) GetStreamingResponse(
	ctx context.Context,
	chatMessages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
) <-chan chatcompletion.ChatStreamingResponse {
	return client.InnerClient.GetStreamingResponse(ctx, chatMessages, options)
}

func (client *CachingChatClient) GetCacheKey(boxed bool, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) string {
	panic("GetCacheKey must be implemented by subclass")
}

func (client *CachingChatClient) ReadCacheAsync(ctx context.Context, key string) (*chatcompletion.ChatResponse, error) {
	panic("ReadCacheAsync must be implemented by subclass")
}

func (client *CachingChatClient) ReadCacheStreamingAsync(ctx context.Context, key string) ([]chatcompletion.ChatResponseUpdate, error) {
	panic("ReadCacheStreamingAsync must be implemented by subclass")
}

func (client *CachingChatClient) WriteCacheAsync(ctx context.Context, key string, value chatcompletion.ChatResponse) error {
	panic("WriteCacheAsync must be implemented by subclass")
}
func (client *CachingChatClient) WriteCacheStreamingAsync(ctx context.Context, key string, value []chatcompletion.ChatResponseUpdate) error {
	panic("WriteCacheStreamingAsync must be implemented by subclass")
}
