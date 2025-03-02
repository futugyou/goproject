package chatcompletion

import "context"

type DelegatingChatClient struct {
	InnerClient IChatClient
}

func NewDelegatingChatClient(delegate IChatClient) *DelegatingChatClient {
	return &DelegatingChatClient{InnerClient: delegate}
}

func (c *DelegatingChatClient) GetResponse(ctx context.Context, chatMessages []ChatMessage, options *ChatOptions) (*ChatResponse, error) {
	return c.InnerClient.GetResponse(ctx, chatMessages, options)
}

func (c *DelegatingChatClient) GetStreamingResponse(ctx context.Context, chatMessages []ChatMessage, options *ChatOptions) <-chan ChatStreamingResponse {
	return c.InnerClient.GetStreamingResponse(ctx, chatMessages, options)
}
