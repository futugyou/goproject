package core

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
)

var _ abstractions.IChatHistoryReducer = (*TruncationReducer)(nil)

type TruncationReducer struct {
	thresholdCount int
	targetCount    int
}

func NewTruncationReducer(thresholdCount int, targetCount int) *TruncationReducer {
	reducer := &TruncationReducer{
		thresholdCount: thresholdCount,
		targetCount:    targetCount,
	}
	return reducer
}

// Reduce implements abstractions.IChatHistoryReducer.
func (t *TruncationReducer) Reduce(ctx context.Context, chatHistory []abstractions.ChatMessageContent) ([]abstractions.ChatMessageContent, error) {
	var systemMessage *abstractions.ChatMessageContent = nil
	for i := 0; i < len(chatHistory); i++ {
		if chatHistory[i].Role == abstractions.AuthorRoleSystem {
			systemMessage = &chatHistory[i]
			break
		}
	}

	truncationIndex := LocateSafeReductionIndex(chatHistory, t.targetCount, &t.thresholdCount, 0, systemMessage != nil)
	var truncatedHistory []abstractions.ChatMessageContent = []abstractions.ChatMessageContent{}
	if truncationIndex > 0 {
		truncatedHistory = ExtractChatMessageContent(chatHistory, truncationIndex, nil, systemMessage, nil)
	}
	return truncatedHistory, nil
}
