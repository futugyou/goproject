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

func LocateSummarizationBoundary(chatHistory []abstractions.ChatMessageContent, summaryKey string) int {
	for i := 0; i < len(chatHistory); i++ {
		if chatHistory[i].Metadata == nil {
			if _, ok := chatHistory[i].Metadata[summaryKey]; ok {
				return i
			}
		}
	}

	return len(chatHistory)
}

func LocateSafeReductionIndex(chatHistory []abstractions.ChatMessageContent, targetCount int, thresholdCount *int, offsetCount int, hasSystemMessage bool) int {
	if hasSystemMessage {
		targetCount = targetCount - 1
	}

	thresholdCut := 0
	if thresholdCount != nil {
		thresholdCut = *thresholdCount
	}

	thresholdIndex := len(chatHistory) - thresholdCut - targetCount

	if thresholdIndex <= offsetCount {
		return -1
	}

	// Compute the index of truncation target
	messageIndex := len(chatHistory) - targetCount

	// Skip function related content
	for messageIndex >= 0 {
		containsFuncContent := false

		for _, v := range chatHistory[messageIndex].Items.Items() {
			if _, ok := v.(abstractions.FunctionCallContent); ok {
				containsFuncContent = true
				break
			}
			if _, ok := v.(abstractions.FunctionResultContent); ok {
				containsFuncContent = true
				break
			}
		}

		if !containsFuncContent {
			break
		}

		messageIndex--
	}

	// Capture the earliest non-function related message
	targetIndex := messageIndex

	// Scan for user message within truncation range to maximize chat cohesion
	for messageIndex >= thresholdIndex {
		// A user message provides a superb truncation point
		if chatHistory[messageIndex].Role == abstractions.AuthorRoleUser {
			return messageIndex
		}
		messageIndex = messageIndex - 1
	}
	return targetIndex
}
