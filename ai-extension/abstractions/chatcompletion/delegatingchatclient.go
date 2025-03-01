package chatcompletion

import "context"

type DelegatingChatClient struct {
	delegate IChatClient
}

func NewDelegatingChatClient(delegate IChatClient) *DelegatingChatClient {
	return &DelegatingChatClient{delegate: delegate}
}

func (c *DelegatingChatClient) GetResponse(ctx context.Context, chatMessages []ChatMessage, options *ChatOptions) (*ChatResponse, error) {
	return c.delegate.GetResponse(ctx, chatMessages, options)
}

func (c *DelegatingChatClient) GetStreamingResponse(ctx context.Context, chatMessages []ChatMessage, options *ChatOptions) <-chan ChatStreamingResponse {
	return c.delegate.GetStreamingResponse(ctx, chatMessages, options)
}
