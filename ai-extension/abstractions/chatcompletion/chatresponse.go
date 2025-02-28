package chatcompletion

import (
	"time"

	"github.com/futugyou/ai-extension/abstractions"
)

type ChatResponse struct {
	Choices              []ChatMessage              `json:"choices"`
	Message              ChatMessage                `json:"-"`
	ResponseId           *string                    `json:"responseId"`
	ChatThreadId         *string                    `json:"chatThreadId"`
	ModelId              *string                    `json:"modelId"`
	CreatedAt            *time.Time                 `json:"createdAt"`
	FinishReason         *ChatFinishReason          `json:"finishReason"`
	Usage                *abstractions.UsageDetails `json:"usage"`
	AdditionalProperties map[string]interface{}     `json:"additionalProperties,omitempty"`
}

type ChatFinishReason string

const (
	ReasonStop          ChatFinishReason = "stop"
	ReasonLength        ChatFinishReason = "length"
	ReasonToolCalls     ChatFinishReason = "tool_calls"
	ReasonContentFilter ChatFinishReason = "content_filter"
)

type ChatResponseUpdate struct {
	AuthorName           *string                `json:"authorName"`
	Role                 *ChatRole              `json:"role"`
	ChoiceIndex          int                    `json:"choiceIndex"`
	Text                 *string                `json:"text"`
	ResponseId           *string                `json:"responseId"`
	ChatThreadId         *string                `json:"chatThreadId"`
	ModelId              *string                `json:"modelId"`
	CreatedAt            *time.Time             `json:"createdAt"`
	FinishReason         *ChatFinishReason      `json:"finishReason"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
	Contents             []interface{}          `json:"contents"`
}

type ChatStreamingResponse struct {
	Update *ChatResponseUpdate
	Err    error
}
