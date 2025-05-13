package chatcompletion

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
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
	var cacheKey = client.GetCacheKey(false, chatMessages, options)
	if cachedResponse, err := client.ReadCache(ctx, cacheKey); err == nil {
		return cachedResponse, nil
	}

	response, err := client.InnerClient.GetResponse(ctx, chatMessages, options)
	if err != nil {
		return nil, err
	}

	client.WriteCache(ctx, cacheKey, *response)
	return response, nil
}

func (client *CachingChatClient) GetStreamingResponse(
	ctx context.Context,
	chatMessages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
) <-chan chatcompletion.ChatStreamingResponse {
	var cacheKey = client.GetCacheKey(true, chatMessages, options)

	if client.CoalesceStreamingUpdates {
		if cachedResponse, err := client.ReadCache(ctx, cacheKey); err == nil {
			streamResp := make(chan chatcompletion.ChatStreamingResponse)

			go func() {
				defer close(streamResp)
				for _, item := range cachedResponse.ToChatResponseUpdates() {
					streamResp <- chatcompletion.ChatStreamingResponse{
						Update: &item,
						Err:    nil,
					}
				}
			}()
			return streamResp
		}

		originalResponse := client.InnerClient.GetStreamingResponse(ctx, chatMessages, options)
		streamResp := make(chan chatcompletion.ChatStreamingResponse)

		go func() {
			defer close(streamResp)
			updates := []chatcompletion.ChatResponseUpdate{}
			for msg := range originalResponse {
				updates = append(updates, *msg.Update)
				streamResp <- msg
			}

			client.WriteCache(ctx, cacheKey, chatcompletion.ToChatResponse(updates))
		}()
		return streamResp
	}

	if cachedResponse, err := client.ReadCacheStreaming(ctx, cacheKey); err == nil {
		streamResp := make(chan chatcompletion.ChatStreamingResponse)

		go func() {
			defer close(streamResp)
			for _, item := range cachedResponse {
				streamResp <- chatcompletion.ChatStreamingResponse{
					Update: &item,
					Err:    nil,
				}
			}
		}()
		return streamResp
	}

	originalResponse := client.InnerClient.GetStreamingResponse(ctx, chatMessages, options)
	streamResp := make(chan chatcompletion.ChatStreamingResponse)

	go func() {
		defer close(streamResp)
		updates := []chatcompletion.ChatResponseUpdate{}
		for msg := range originalResponse {
			updates = append(updates, *msg.Update)
			streamResp <- msg
		}

		client.WriteCacheStreaming(ctx, cacheKey, updates)
	}()
	return streamResp
}

func (client *CachingChatClient) GetCacheKey(values ...interface{}) string {
	panic("GetCacheKey must be implemented by subclass")
}

func (client *CachingChatClient) ReadCache(ctx context.Context, key string) (*chatcompletion.ChatResponse, error) {
	panic("ReadCache must be implemented by subclass")
}

func (client *CachingChatClient) ReadCacheStreaming(ctx context.Context, key string) ([]chatcompletion.ChatResponseUpdate, error) {
	panic("ReadCacheStreaming must be implemented by subclass")
}

func (client *CachingChatClient) WriteCache(ctx context.Context, key string, value chatcompletion.ChatResponse) error {
	panic("WriteCache must be implemented by subclass")
}
func (client *CachingChatClient) WriteCacheStreaming(ctx context.Context, key string, value []chatcompletion.ChatResponseUpdate) error {
	panic("WriteCacheStreaming must be implemented by subclass")
}
