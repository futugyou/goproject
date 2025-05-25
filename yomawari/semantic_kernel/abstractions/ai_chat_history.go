package abstractions

import (
	"encoding/base64"

	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"golang.org/x/net/context"
)

type ChatHistory struct {
	core.List[ChatMessageContent]
}

func NewChatHistory(message string, role AuthorRole) *ChatHistory {
	c := &ChatHistory{
		List: *core.NewList[ChatMessageContent](),
	}
	c.Add(ChatMessageContent{
		Role:    role,
		Content: message,
	})
	return c
}

func NewSystemChatHistory(message string) *ChatHistory {
	c := &ChatHistory{
		List: *core.NewList[ChatMessageContent](),
	}
	c.Add(ChatMessageContent{
		Role:    AuthorRoleSystem,
		Content: message,
	})
	return c
}

func NewChatHistoryWithContents(messages []ChatMessageContent) *ChatHistory {
	c := &ChatHistory{
		List: *core.NewList[ChatMessageContent](),
	}
	c.AddRange(messages)
	return c
}

func (c *ChatHistory) AddMessage(authorRole AuthorRole, content string, encoding *base64.Encoding, metadata map[string]any, colls *ChatMessageContentItemCollection) {
	c.List.Add(ChatMessageContent{
		Metadata: metadata,
		Role:     authorRole,
		Items:    colls,
		Content:  content,
		Encoding: encoding,
	})
}

func (c *ChatHistory) AddUserMessage(content string, colls *ChatMessageContentItemCollection) {
	c.List.Add(ChatMessageContent{
		Role:    AuthorRoleUser,
		Items:   colls,
		Content: content,
	})
}

func (c *ChatHistory) AddAssistantMessage(content string) {
	c.List.Add(ChatMessageContent{
		Role:    AuthorRoleAssistant,
		Content: content,
	})
}

func (c *ChatHistory) AddSystemMessage(content string) {
	c.List.Add(ChatMessageContent{
		Role:    AuthorRoleSystem,
		Content: content,
	})
}

func (c *ChatHistory) AddDeveloperMessage(content string) {
	c.List.Add(ChatMessageContent{
		Role:    AuthorRoleDeveloper,
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

func (chatHistory *ChatHistory) ToChatMessageList() []chatcompletion.ChatMessage {
	result := []chatcompletion.ChatMessage{}
	for _, v := range chatHistory.Items() {
		result = append(result, v.ToChatMessage())
	}
	return result
}
