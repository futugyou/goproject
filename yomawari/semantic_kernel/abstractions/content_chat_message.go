package abstractions

import (
	"encoding/base64"
	"encoding/json"

	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	aicontents "github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
)

type AuthorRole string

const (
	AuthorRoleDeveloper AuthorRole = "developer"
	AuthorRoleSystem    AuthorRole = "system"
	AuthorRoleAssistant AuthorRole = "assistant"
	AuthorRoleUser      AuthorRole = "user"
	AuthorRoleTool      AuthorRole = "tool"
)

func CreateAuthorRole(name string) AuthorRole {
	switch name {
	case "developer":
		return AuthorRoleDeveloper
	case "system":
		return AuthorRoleSystem
	case "assistant":
		return AuthorRoleAssistant
	case "user":
		return AuthorRoleUser
	case "tool":
		return AuthorRoleTool
	default:
		return AuthorRoleUser
	}
}

type ChatMessageContentItemCollection struct {
	core.List[KernelContent] `json:"items"`
}

func (c ChatMessageContentItemCollection) MarshalJSON() ([]byte, error) {
	var rawItems []json.RawMessage
	for _, item := range c.Items() {
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
		c.Add(content)
	}
	return nil
}

type ChatMessageContent struct {
	MimeType     string                            `json:"mimeType"`
	ModelId      string                            `json:"modelId"`
	Metadata     map[string]any                    `json:"metadata"`
	AuthorName   string                            `json:"authorName"`
	Role         AuthorRole                        `json:"role"`
	Items        *ChatMessageContentItemCollection `json:"items"`
	Content      string                            `json:"-"`
	Encoding     *base64.Encoding                  `json:"-"`
	Source       any                               `json:"-"`
	InnerContent any                               `json:"-"`
}

func (ChatMessageContent) Type() string {
	return "chatMessage"
}

func (f ChatMessageContent) Hash() string {
	return f.ToString()
}

func (f ChatMessageContent) ToString() string {
	return f.Content
}

func (f ChatMessageContent) GetInnerContent() any {
	return f.InnerContent
}

func (c *ChatMessageContent) GetFunctionCalls() []FunctionCallContent {
	var result []FunctionCallContent
	for _, item := range c.Items.Items() {
		if item.Type() == "functionCall" {
			result = append(result, item.(FunctionCallContent))
		}
	}
	return result
}

func (c ChatMessageContent) GetContent() string {
	for _, item := range c.Items.Items() {
		if textContent, ok := item.(TextContent); ok && item.Type() == "streaming-function-call-update" {
			return textContent.Text
		}
	}
	return ""
}

func (c *ChatMessageContent) SetContent(content string) {
	for i, item := range c.Items.Items() {
		if textContent, ok := item.(TextContent); ok && item.Type() == "streaming-function-call-update" {
			textContent.Text = content
			c.Items.Set(i, textContent)
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
	c.Items.Add(textContent)
}

func (content *ChatMessageContent) ToChatMessage() chatcompletion.ChatMessage {
	message := &chatcompletion.ChatMessage{
		AdditionalProperties: content.Metadata,
		AuthorName:           &content.AuthorName,
		RawRepresentation:    content.InnerContent,
	}

	message.Role = chatcompletion.StringToChatRole(string(content.Role))

	for _, item := range content.Items.Items() {
		var aiContent aicontents.IAIContent

		switch tc := item.(type) {
		case TextContent:
			aiContent = aicontents.NewTextContent(tc.Text)
		case ImageContent:
			if len(tc.DataUri) > 0 {
				aiContent = aicontents.NewDataContent(tc.DataUri, tc.MimeType)
			}

			if len(tc.Uri.String()) > 0 {
				mimeType := "image/*"
				if len(tc.MimeType) > 0 {
					mimeType = tc.MimeType
				}
				aiContent = aicontents.UriContent{URI: tc.Uri.String(), MediaType: mimeType}
			}

		case AudioContent:
			if len(tc.DataUri) > 0 {
				aiContent = aicontents.NewDataContent(tc.DataUri, tc.MimeType)
			}

			if len(tc.Uri.String()) > 0 {
				mimeType := "audio/*"
				if len(tc.MimeType) > 0 {
					mimeType = tc.MimeType
				}
				aiContent = aicontents.UriContent{URI: tc.Uri.String(), MediaType: mimeType}
			}
		case BinaryContent:
			if len(tc.DataUri) > 0 {
				aiContent = aicontents.NewDataContent(tc.DataUri, tc.MimeType)
			}

			if len(tc.Uri.String()) > 0 {
				mimeType := "application/octet-stream"
				if len(tc.MimeType) > 0 {
					mimeType = tc.MimeType
				}
				aiContent = aicontents.UriContent{URI: tc.Uri.String(), MediaType: mimeType}
			}
		case FunctionCallContent:
			aiContent = aicontents.FunctionCallContent{CallId: tc.Id, Name: tc.FunctionName, Arguments: tc.Arguments}
		case FunctionResultContent:
			aiContent = aicontents.FunctionResultContent{CallId: tc.CallId, Result: tc.Result}
		}

		if aiContent != nil {
			message.Contents = append(message.Contents, aiContent)
		}
	}

	return *message
}
