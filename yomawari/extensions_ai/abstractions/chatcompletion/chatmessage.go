package chatcompletion

import (
	"encoding/json"
	"fmt"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
)

type ChatMessage struct {
	AuthorName           *string                `json:"authorName"`
	Role                 ChatRole               `json:"role"`
	Contents             []contents.IAIContent  `json:"contents"`
	MessageId            *string                `json:"messageId"`
	RawRepresentation    interface{}            `json:"-"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
}

func NewChatMessage(role ChatRole, contents []contents.IAIContent) *ChatMessage {
	return &ChatMessage{
		Role:     role,
		Contents: contents,
	}
}

func NewChatMessageWithText(role ChatRole, text string) *ChatMessage {
	return &ChatMessage{
		Role: role,
		Contents: []contents.IAIContent{
			contents.NewTextContent(text),
		},
	}
}

func (c *ChatMessage) Text() string {
	return contents.ConcatTextContents(c.Contents)
}

func (cru *ChatMessage) UnmarshalJSON(data []byte) error {
	temp := struct {
		Role                 ChatRole               `json:"role"`
		MessageId            *string                `json:"messageId"`
		AuthorName           *string                `json:"authorName"`
		AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
		Contents             []json.RawMessage      `json:"contents"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	cru.Role = temp.Role
	cru.MessageId = temp.MessageId
	cru.AuthorName = temp.AuthorName
	cru.AdditionalProperties = temp.AdditionalProperties
	cru.RawRepresentation = json.RawMessage(data)

	for _, raw := range temp.Contents {
		var base struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &base); err != nil {
			return err
		}

		content, err := createContentByType(base.Type)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(raw, content); err != nil {
			return err
		}

		cru.Contents = append(cru.Contents, content)
	}

	return nil
}

func (cru *ChatMessage) MarshalJSON() ([]byte, error) {
	var contentsRaw []json.RawMessage
	for _, content := range cru.Contents {
		raw, err := json.Marshal(content)
		if err != nil {
			return nil, err
		}
		contentsRaw = append(contentsRaw, raw)
	}

	temp := struct {
		Role                 ChatRole               `json:"role"`
		MessageId            *string                `json:"messageId"`
		AuthorName           *string                `json:"authorName"`
		AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
		Contents             []json.RawMessage      `json:"contents"`
	}{
		Role:                 cru.Role,
		MessageId:            cru.MessageId,
		AuthorName:           cru.AuthorName,
		AdditionalProperties: cru.AdditionalProperties,
		Contents:             contentsRaw,
	}

	return json.Marshal(temp)
}

func createContentByType(typeStr string) (contents.IAIContent, error) {
	if factory, ok := contents.ContentTypeRegistry[typeStr]; ok {
		return factory(), nil
	}
	return nil, fmt.Errorf("unknown type: %s", typeStr)
}

type ChatRole string

const (
	RoleSystem    ChatRole = "system"
	RoleAssistant ChatRole = "assistant"
	RoleUser      ChatRole = "user"
	RoleTool      ChatRole = "tool"
	RoleDeveloper ChatRole = "developer"
)

func StringToChatRole(s string) ChatRole {
	switch s {
	case string(RoleSystem):
		return RoleSystem
	case string(RoleAssistant):
		return RoleAssistant
	case string(RoleUser):
		return RoleUser
	case string(RoleDeveloper):
		return RoleDeveloper
	case string(RoleTool):
		return RoleTool
	default:
		return RoleSystem
	}
}
