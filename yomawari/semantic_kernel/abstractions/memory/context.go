package memory

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
)

type AIContext struct {
	Instructions string
	AIFunctions  []functions.AIFunction
}

type AIContextProvider interface {
	ConversationCreated(ctx context.Context, conversationId string) error
	MessageAdding(ctx context.Context, conversationId string, newMessage chatcompletion.ChatMessage) error
	ConversationDeleting(ctx context.Context, conversationId string) error
	ModelInvoking(ctx context.Context, newMessages []chatcompletion.ChatMessage) (*AIContext, error)
	Suspending(ctx context.Context, conversationId string) error
	Resuming(ctx context.Context, conversationId string) error
}
