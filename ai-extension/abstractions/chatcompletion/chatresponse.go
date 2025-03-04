package chatcompletion

import (
	"encoding/json"
	"fmt"
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

func (c *ChatResponse) ToChatResponseUpdates() []ChatResponseUpdate {
	var updates []ChatResponseUpdate

	for i, choice := range c.Choices {
		update := ChatResponseUpdate{
			ChatThreadId:         c.ChatThreadId,
			ChoiceIndex:          i,
			AdditionalProperties: choice.AdditionalProperties,
			AuthorName:           choice.AuthorName,
			Contents:             choice.Contents,
			Role:                 &choice.Role,
			ResponseId:           c.ResponseId,
			CreatedAt:            c.CreatedAt,
			FinishReason:         c.FinishReason,
			ModelId:              c.ModelId,
		}
		updates = append(updates, update)
	}

	if c.AdditionalProperties != nil || c.Usage != nil {
		extra := ChatResponseUpdate{
			AdditionalProperties: c.AdditionalProperties,
		}

		if c.Usage != nil {
			extra.Contents = append(extra.Contents, contents.UsageContent{
				AIContent: contents.AIContent{},
				Details:   *c.Usage,
			})
		}

		updates = append(updates, extra)
	}

	return updates
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
	Contents             []contents.IAIContent  `json:"contents"`
}

func (cru *ChatResponseUpdate) UnmarshalJSON(data []byte) error {
	temp := struct {
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
		Contents             []json.RawMessage      `json:"contents"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	for _, raw := range temp.Contents {
		var base struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &base); err != nil {
			return err
		}

		var content contents.IAIContent
		switch base.Type {
		case "AIContent":
			content = &contents.AIContent{}
		case "FunctionCallContent":
			content = &contents.FunctionCallContent{}
		case "FunctionResultContent":
			content = &contents.FunctionResultContent{}
		case "TextContent":
			content = &contents.TextContent{}
		case "UsageContent":
			content = &contents.UsageContent{}
		default:
			return fmt.Errorf("unknown type: %s", base.Type)
		}

		if err := json.Unmarshal(raw, content); err != nil {
			return err
		}

		cru.Contents = append(cru.Contents, content)
	}
	return nil
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
			Contents:             []contents.IAIContent{},
			AuthorName:           new(string),
			AdditionalProperties: map[string]interface{}{},
		}
		if v, ok := messages[update.ChoiceIndex]; ok {
			message = v
		}

		for _, con := range update.Contents {
			switch c := con.(type) {
			case contents.UsageContent:
				response.Usage.AddUsageDetails(c.Details)
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

func coalesceTextContent(conts []contents.IAIContent) []contents.IAIContent {
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

func filterNonNilContents(conts []contents.IAIContent) []contents.IAIContent {
	result := []contents.IAIContent{}
	for _, content := range conts {
		if content != nil {
			result = append(result, content)
		}
	}
	return result
}
