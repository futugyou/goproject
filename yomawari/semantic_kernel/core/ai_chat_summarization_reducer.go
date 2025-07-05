package core

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
)

const SummaryMetadataKey = "__summary__"
const DefaultSummarizationPrompt = `
Provide a concise and complete summarization of the entire dialog that does not exceed 5 sentences

This summary must always:
- Consider both user and assistant interactions
- Maintain continuity for the purpose of further dialog
- Include details from any existing summary
- Focus on the most significant aspects of the dialog

This summary must never:
- Critique, correct, interpret, presume, or assume
- Identify faults, mistakes, misunderstanding, or correctness
- Analyze what has not occurred
- Exclude details from any existing summary
`

type SummarizationReducer struct {
	SummarizationInstructions string
	FailOnError               bool
	UseSingleSummary          bool
	service                   abstractions.IChatCompletionService
	targetCount               int
	thresholdCount            int
}

func NewSummarizationReducer(service abstractions.IChatCompletionService, targetCount int, thresholdCount *int) *SummarizationReducer {
	reducer := &SummarizationReducer{
		SummarizationInstructions: DefaultSummarizationPrompt,
		FailOnError:               true,
		UseSingleSummary:          true,
		service:                   service,
		targetCount:               targetCount,
		thresholdCount:            0,
	}

	if thresholdCount != nil {
		reducer.thresholdCount = *thresholdCount
	}
	return reducer
}

// Reduce implements abstractions.IChatHistoryReducer.
func (s *SummarizationReducer) Reduce(ctx context.Context, chatHistory []abstractions.ChatMessageContent) ([]abstractions.ChatMessageContent, error) {
	var systemMessage *abstractions.ChatMessageContent = nil
	for i := 0; i < len(chatHistory); i++ {
		if chatHistory[i].Role == abstractions.AuthorRoleSystem {
			systemMessage = &chatHistory[i]
			break
		}
	}

	insertionPoint := LocateSummarizationBoundary(chatHistory, SummaryMetadataKey)
	truncationIndex := LocateSafeReductionIndex(chatHistory, s.targetCount, &s.thresholdCount, insertionPoint, systemMessage != nil)

	var truncatedHistory []abstractions.ChatMessageContent = []abstractions.ChatMessageContent{}

	assemblySummarizedHistory := func(summaryMessage *abstractions.ChatMessageContent, systemMessage *abstractions.ChatMessageContent) []abstractions.ChatMessageContent {
		result := []abstractions.ChatMessageContent{}
		if systemMessage == nil {
			return result
		}
		if insertionPoint > 0 && !s.UseSingleSummary {
			for i := 0; i < insertionPoint; i++ {
				result = append(result, chatHistory[i])
			}
		}
		if summaryMessage != nil {
			result = append(result, *summaryMessage)
		}
		for i := truncationIndex; i < len(chatHistory); i++ {
			result = append(result, chatHistory[i])
		}
		return result
	}

	if truncationIndex >= 0 {
		startIndex := insertionPoint
		if s.UseSingleSummary {
			startIndex = 0
		}
		summarizedHistory := ExtractChatMessageContent(chatHistory, startIndex, &truncationIndex, nil, func(chatMessageContent abstractions.ChatMessageContent) bool {
			if chatMessageContent.Items == nil {
				return false
			}

			for _, v := range chatMessageContent.Items.Items() {
				if _, ok := v.(abstractions.FunctionCallContent); ok {
					return false
				}
				if _, ok := v.(abstractions.FunctionResultContent); ok {
					return false
				}
			}
			return false
		})
		systemMessage := abstractions.ChatMessageContent{
			Role:    abstractions.AuthorRoleSystem,
			Content: s.SummarizationInstructions,
		}
		summarizedHistory = append(summarizedHistory, systemMessage)
		summarizationRequest := abstractions.NewChatHistoryWithContents(summarizedHistory)
		summaryMessages, err := s.service.GetChatMessageContents(ctx, *summarizationRequest, nil, nil)
		if err != nil || len(summaryMessages) == 0 {
			return truncatedHistory, err
		}
		summaryMessage := summaryMessages[0]
		summaryMessage.Metadata = map[string]interface{}{SummaryMetadataKey: true}
		truncatedHistory = assemblySummarizedHistory(&summaryMessage, &systemMessage)
	}

	return truncatedHistory, nil
}

var _ abstractions.IChatHistoryReducer = (*SummarizationReducer)(nil)
