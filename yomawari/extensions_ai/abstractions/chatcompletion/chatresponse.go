package chatcompletion

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/futugyou/yomawari/extensions_ai/abstractions"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
)

type ChatResponse struct {
	Messages   []ChatMessage `json:"choices"`
	ResponseId *string       `json:"responseId"`
	//Obsolete use ConversationId instead
	ChatThreadId         *string                    `json:"chatThreadId"`
	ConversationId       *string                    `json:"conversationId,omitempty"`
	ModelId              *string                    `json:"modelId"`
	CreatedAt            *time.Time                 `json:"createdAt"`
	FinishReason         *ChatFinishReason          `json:"finishReason"`
	Usage                *abstractions.UsageDetails `json:"usage"`
	RawRepresentation    interface{}                `json:"-"`
	AdditionalProperties map[string]interface{}     `json:"additionalProperties,omitempty"`
}

func NewChatResponse(messages []ChatMessage, message *ChatMessage) *ChatResponse {
	response := &ChatResponse{
		Messages: []ChatMessage{},
	}
	if len(messages) > 0 {
		response.Messages = messages
	}
	if message != nil {
		response.Messages = append(response.Messages, *message)
	}
	return response
}

func (c *ChatResponse) Text() string {
	return ConcatMessagesContents(c.Messages)
}

func ConcatMessagesContents(contents []ChatMessage) string {
	var text string
	for _, content := range contents {
		text += content.Text()
	}
	return text
}

func (c *ChatResponse) ToChatResponseUpdates() []ChatResponseUpdate {
	var updates []ChatResponseUpdate

	for _, choice := range c.Messages {
		update := ChatResponseUpdate{
			ChatThreadId:         c.ChatThreadId,
			ConversationId:       c.ConversationId,
			AdditionalProperties: choice.AdditionalProperties,
			AuthorName:           choice.AuthorName,
			Contents:             choice.Contents,
			MessageId:            choice.MessageId,
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
			extra.Contents = append(extra.Contents, contents.NewUsageContent(*c.Usage))
		}

		updates = append(updates, extra)
	}

	return updates
}

type ChatFinishReason string

const (
	ReasonStop          = ChatFinishReason("stop")
	ReasonLength        = ChatFinishReason("length")
	ReasonToolCalls     = ChatFinishReason("tool_calls")
	ReasonContentFilter = ChatFinishReason("content_filter")
	ReasonUnknown       = ChatFinishReason("unknown")
)

type ChatResponseUpdate struct {
	AuthorName           *string                `json:"authorName"`
	Role                 *ChatRole              `json:"role"`
	Text                 *string                `json:"-"`
	Contents             []contents.IAIContent  `json:"contents"`
	RawRepresentation    interface{}            `json:"-"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
	ResponseId           *string                `json:"responseId"`
	MessageId            *string                `json:"messageId"`
	//Obsolete use ConversationId instead
	ChatThreadId   *string           `json:"chatThreadId"`
	ConversationId *string           `json:"conversationId,omitempty"`
	CreatedAt      *time.Time        `json:"createdAt"`
	FinishReason   *ChatFinishReason `json:"finishReason"`
	ModelId        *string           `json:"modelId"`
}

func (cru *ChatResponseUpdate) UnmarshalJSON(data []byte) error {
	type Alias ChatResponseUpdate
	temp := &struct {
		*Alias
		Contents []json.RawMessage `json:"contents"`
	}{
		Alias: (*Alias)(cru),
	}

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

func ToChatResponse(updates []ChatResponseUpdate) ChatResponse {
	if len(updates) == 0 {
		return ChatResponse{}
	}

	response := new(ChatResponse)
	for _, update := range updates {
		response.ChatThreadId = update.ChatThreadId
		response.ConversationId = update.ConversationId
		response.CreatedAt = update.CreatedAt
		response.FinishReason = update.FinishReason
		response.ModelId = update.ModelId
		response.ResponseId = update.ResponseId

		var message ChatMessage
		if len(response.Messages) == 0 || response.ResponseId != nil && update.ResponseId != nil && *response.ResponseId != *update.ResponseId {
			message = ChatMessage{
				AuthorName:           new(string),
				Role:                 RoleAssistant,
				Contents:             []contents.IAIContent{},
				AdditionalProperties: map[string]interface{}{},
			}
		} else {
			message = response.Messages[len(response.Messages)-1]
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

		response.Messages = append(response.Messages, message)
	}

	finalizeResponse(response)
	return *response
}

func finalizeResponse(response *ChatResponse) {
	for i := 0; i < len(response.Messages); i++ {
		response.Messages[i].Contents = CoalesceTextContent(response.Messages[i].Contents)
	}
}

func CoalesceTextContent(conts []contents.IAIContent) []contents.IAIContent {
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
			AIContent: contents.NewAIContent(nil, nil),
			Text:      coalescedText.String(),
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
