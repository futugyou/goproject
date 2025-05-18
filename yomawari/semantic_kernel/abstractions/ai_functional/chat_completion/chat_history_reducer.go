package chat_completion

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
)

type IChatHistoryReducer interface {
	Reduce(ctx context.Context, chatHistory []contents.ChatMessageContent) ([]contents.ChatMessageContent, error)
}
