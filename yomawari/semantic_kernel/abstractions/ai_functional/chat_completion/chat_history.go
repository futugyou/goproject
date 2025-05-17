package chat_completion

import (
	"encoding/base64"

	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
)

type ChatHistory struct {
	core.List[contents.ChatMessageContent]
}

func NewChatHistory(message string, role contents.AuthorRole) *ChatHistory {
	c := &ChatHistory{
		List: *core.NewList[contents.ChatMessageContent](),
	}
	c.Add(contents.ChatMessageContent{
		Role:    role,
		Content: message,
	})
	return c
}

func NewSystemChatHistory(message string) *ChatHistory {
	c := &ChatHistory{
		List: *core.NewList[contents.ChatMessageContent](),
	}
	c.Add(contents.ChatMessageContent{
		Role:    contents.AuthorRoleSystem,
		Content: message,
	})
	return c
}

func NewChatHistoryWithContents(messages []contents.ChatMessageContent) *ChatHistory {
	c := &ChatHistory{
		List: *core.NewList[contents.ChatMessageContent](),
	}
	c.AddRange(messages)
	return c
}

func (c *ChatHistory) AddMessage(authorRole contents.AuthorRole, content string, encoding *base64.Encoding, metadata map[string]any, colls *contents.ChatMessageContentItemCollection) {
	c.List.Add(contents.ChatMessageContent{
		Metadata: metadata,
		Role:     authorRole,
		Items:    colls,
		Content:  content,
		Encoding: encoding,
	})
}

func (c *ChatHistory) AddUserMessage(content string, colls *contents.ChatMessageContentItemCollection) {
	c.List.Add(contents.ChatMessageContent{
		Role:    contents.AuthorRoleUser,
		Items:   colls,
		Content: content,
	})
}

func (c *ChatHistory) AddAssistantMessage(content string) {
	c.List.Add(contents.ChatMessageContent{
		Role:    contents.AuthorRoleAssistant,
		Content: content,
	})
}

func (c *ChatHistory) AddSystemMessage(content string) {
	c.List.Add(contents.ChatMessageContent{
		Role:    contents.AuthorRoleSystem,
		Content: content,
	})
}

func (c *ChatHistory) AddDeveloperMessage(content string) {
	c.List.Add(contents.ChatMessageContent{
		Role:    contents.AuthorRoleDeveloper,
		Content: content,
	})
}
