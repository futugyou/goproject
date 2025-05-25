package abstractions

import (
	"context"
)

type IChatHistoryReducer interface {
	Reduce(ctx context.Context, chatHistory []ChatMessageContent) ([]ChatMessageContent, error)
}
