package chatcompletion

import (
	"sort"
	"strings"
	"time"

	"github.com/futugyou/ai-extension/abstractions"
	"github.com/futugyou/ai-extension/abstractions/contents"
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

func ToChatResponse(updates []ChatResponseUpdate, coalesceContent bool) ChatResponse {
	if len(updates) == 0 {
		return ChatResponse{}
	}

	response := new(ChatResponse)
	messages := map[int]ChatMessage{}

	for _, update := range updates {
		response.ChatThreadId = update.ChatThreadId
		response.CreatedAt = update.CreatedAt
		response.FinishReason = update.FinishReason
		response.ModelId = update.ModelId
		response.ResponseId = update.ResponseId

		message := ChatMessage{
			Role:                 "",
			Message:              "",
			Text:                 new(string),
			Contents:             []interface{}{},
			AuthorName:           new(string),
			AdditionalProperties: map[string]interface{}{},
		}
		if v, ok := messages[update.ChoiceIndex]; ok {
			message = v
		}

		for _, con := range update.Contents {
			switch c := con.(type) {
			case contents.UsageContent:
				response.Usage.Add(c.Details)
			default:
				message.Contents = append(message.Contents, c)
			}
		}
		message.AuthorName = update.AuthorName

		if update.Role != nil {
			message.Role = *update.Role
		}

		for key, vaule := range update.AdditionalProperties {
			message.AdditionalProperties[key] = vaule

		}
		messages[update.ChoiceIndex] = message
	}

	addMessagesToResponse(messages, response, coalesceContent)
	return *response
}

func addMessagesToResponse(messages map[int]ChatMessage, response *ChatResponse, coalesceContent bool) {
	if len(messages) <= 1 {
		for _, message := range messages {
			addMessage(response, coalesceContent, message)
		}

		if len(response.Choices) == 1 && response.Choices[0].AdditionalProperties != nil {
			response.AdditionalProperties = response.Choices[0].AdditionalProperties
			response.Choices[0].AdditionalProperties = nil
		}
	} else {
		keys := make([]int, 0, len(messages))
		for key := range messages {
			keys = append(keys, key)
		}
		sort.Ints(keys)

		for _, key := range keys {
			addMessage(response, coalesceContent, messages[key])
		}
	}
}

func addMessage(response *ChatResponse, coalesceContent bool, message ChatMessage) {
	if message.Role == "" {
		message.Role = "Assistant"
	}

	if coalesceContent {
		message.Contents = coalesceTextContent(message.Contents)
	}

	response.Choices = append(response.Choices, message)
}

func coalesceTextContent(conts []interface{}) []interface{} {
	var coalescedText *strings.Builder

	start := 0
	for start < len(conts)-1 {
		firstText, ok1 := conts[start].(*contents.TextContent)
		if !ok1 {
			start++
			continue
		}

		secondText, ok2 := conts[start+1].(*contents.TextContent)
		if !ok2 {
			start += 2
			continue
		}

		if coalescedText == nil {
			coalescedText = &strings.Builder{}
		}
		coalescedText.Reset()
		coalescedText.WriteString(firstText.Text)
		coalescedText.WriteString(secondText.Text)
		conts[start+1] = nil

		i := start + 2
		for i < len(conts) {
			next, ok := conts[i].(*contents.TextContent)
			if !ok {
				break
			}
			coalescedText.WriteString(next.Text)
			conts[i] = nil
			i++
		}

		conts[start] = &contents.TextContent{
			AIContent: contents.AIContent{
				AdditionalProperties: firstText.AdditionalProperties,
			},
			Text: coalescedText.String(),
		}

		start = i
	}

	conts = filterNonNilContents(conts)
	return conts
}

func filterNonNilContents(contents []interface{}) []interface{} {
	result := []interface{}{}
	for _, content := range contents {
		if content != nil {
			result = append(result, content)
		}
	}
	return result
}
