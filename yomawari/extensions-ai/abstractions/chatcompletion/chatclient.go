package chatcompletion

import "context"

type IChatClient interface {
	GetResponse(ctx context.Context, chatMessages []ChatMessage, options *ChatOptions) (*ChatResponse, error)
	GetStreamingResponse(ctx context.Context, chatMessages []ChatMessage, options *ChatOptions) <-chan ChatStreamingResponse
}
