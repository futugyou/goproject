package chatcompletion

import (
	"encoding/json"
	"fmt"

	"github.com/futugyou/yomawari/generative-ai/abstractions/contents"
)

type ChatMessage struct {
	AuthorName           *string                `json:"authorName"`
	Role                 ChatRole               `json:"role"`
	Contents             []contents.IAIContent  `json:"contents"`
	RawRepresentation    interface{}            `json:"-"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
}

func (c *ChatMessage) Text() string {
	return contents.ConcatTextContents(c.Contents)
}

func (cru *ChatMessage) UnmarshalJSON(data []byte) error {
	temp := struct {
		Role                 ChatRole               `json:"role"`
		Message              string                 `json:"message"`
		Text                 *string                `json:"-"`
		AuthorName           *string                `json:"authorName"`
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

type ChatRole string

const (
	RoleSystem    ChatRole = "system"
	RoleAssistant ChatRole = "assistant"
	RoleUser      ChatRole = "user"
	RoleTool      ChatRole = "tool"
)

func StringToChatRole(s string) ChatRole {
	switch s {
	case string(RoleSystem):
		return RoleSystem
	case string(RoleAssistant):
		return RoleAssistant
	case string(RoleUser):
		return RoleUser
	case string(RoleTool):
		return RoleTool
	default:
		return RoleSystem
	}
}
