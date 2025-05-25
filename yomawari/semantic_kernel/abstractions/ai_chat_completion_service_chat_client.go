package abstractions

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
)

var _ chatcompletion.IChatClient = (*ChatCompletionServiceChatClient)(nil)

type ChatCompletionServiceChatClient struct {
}

// GetResponse implements chatcompletion.IChatClient.
func (c *ChatCompletionServiceChatClient) GetResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	panic("unimplemented")
}

// GetStreamingResponse implements chatcompletion.IChatClient.
func (c *ChatCompletionServiceChatClient) GetStreamingResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	panic("unimplemented")
}
