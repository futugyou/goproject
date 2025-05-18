package chat_completion

import (
	"encoding/base64"

	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
	"golang.org/x/net/context"
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

func (c *ChatHistory) ReduceInPlace(ctx context.Context, reducer IChatHistoryReducer) bool {
	if reducer == nil {
		return false
	}

	reducedHistory, err := reducer.Reduce(ctx, c.Items())
	if err != nil {
		return false
	}
	if reducedHistory == nil {
		return false
	}

	c.Clear()
	c.AddRange(reducedHistory)

	return true
}

func (c *ChatHistory) Reduce(ctx context.Context, reducer IChatHistoryReducer) (*ChatHistory, error) {
	if reducer != nil {
		reducedHistory, err := reducer.Reduce(ctx, c.Items())
		if err != nil {
			return c, err
		}
		return NewChatHistoryWithContents(reducedHistory), nil
	}

	return c, nil
}
