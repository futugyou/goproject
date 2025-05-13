package contents

import (
	"encoding/base64"
	"encoding/json"
)

type AuthorRole string

const (
	AuthorNameDeveloper AuthorRole = "developer"
	AuthorRoleSystem    AuthorRole = "system"
	AuthorRoleAssistant AuthorRole = "assistant"
	AuthorRoleUser      AuthorRole = "user"
	AuthorRoleTool      AuthorRole = "tool"
)

type ChatMessageContentItemCollection struct {
	Items []KernelContent `json:"items"`
}

func (c ChatMessageContentItemCollection) MarshalJSON() ([]byte, error) {
	var rawItems []json.RawMessage
	for _, item := range c.Items {
		b, err := MarshalKernelContent(item)
		if err != nil {
			return nil, err
		}
		rawItems = append(rawItems, b)
	}
	return json.Marshal(map[string]any{
		"items": rawItems,
	})
}

func (c *ChatMessageContentItemCollection) UnmarshalJSON(data []byte) error {
	var raw struct {
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, item := range raw.Items {
		content, err := UnmarshalKernelContent(item)
		if err != nil {
			return err
		}
		c.Items = append(c.Items, content)
	}
	return nil
}

type ChatMessageContent struct {
	MimeType     string                           `json:"mimeType"`
	ModelId      string                           `json:"modelId"`
	Metadata     map[string]any                   `json:"metadata"`
	AuthorName   string                           `json:"authorName"`
	Role         AuthorRole                       `json:"role"`
	Items        ChatMessageContentItemCollection `json:"items"`
	Content      string                           `json:"-"`
	Encoding     *base64.Encoding                 `json:"-"`
	Source       any                              `json:"-"`
	InnerContent any                              `json:"-"`
}

func (ChatMessageContent) Type() string {
	return "chatMessage"
}

func (c *ChatMessageContent) GetFunctionCalls() []FunctionCallContent {
	var result []FunctionCallContent
	for _, item := range c.Items.Items {
		if item.Type() == "functionCall" {
			result = append(result, item.(FunctionCallContent))
		}
	}
	return result
}

func (c ChatMessageContent) GetContent() string {
	for _, item := range c.Items.Items {
		if textContent, ok := item.(TextContent); ok && item.Type() == "streaming-function-call-update" {
			return textContent.Text
		}
	}
	return ""
}

func (c *ChatMessageContent) SetContent(content string) {
	for i, item := range c.Items.Items {
		if textContent, ok := item.(TextContent); ok && item.Type() == "streaming-function-call-update" {
			textContent.Text = content
			c.Items.Items[i] = textContent
			return
		}
	}

	var textContent TextContent = TextContent{
		ModelId:      c.ModelId,
		Metadata:     c.Metadata,
		InnerContent: c.InnerContent,
		Text:         content,
		Encoding:     c.Encoding,
	}
	c.Items.Items = append(c.Items.Items, textContent)
}
