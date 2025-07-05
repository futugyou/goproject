package core

import (
	"math"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
)

func ExtractChatMessageContent(chatHistory []abstractions.ChatMessageContent, startIndex int, finalIndex *int, systemMessage *abstractions.ChatMessageContent, filter func(chatMessageContent abstractions.ChatMessageContent) bool) []abstractions.ChatMessageContent {
	result := []abstractions.ChatMessageContent{}
	maxIndex := len(chatHistory) - 1
	if startIndex > maxIndex {
		return result
	}

	if systemMessage != nil {
		result = append(result, *systemMessage)
	}

	finalIdx := maxIndex
	if finalIndex != nil {
		finalIdx = int(math.Min((float64)(*finalIndex), (float64)(maxIndex)))
	}

	for i := startIndex; i <= finalIdx; i++ {
		if filter != nil && filter(chatHistory[i]) {
			continue
		}

		result = append(result, chatHistory[i])
	}

	return result
}
