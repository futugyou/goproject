package evaluation

import "github.com/futugyou/yomawari/extensions-ai/abstractions/chatcompletion"

type ChatConfiguration struct {
	ChatClient chatcompletion.IChatClient
}

func NewChatConfiguration(chatClient chatcompletion.IChatClient) *ChatConfiguration {
	return &ChatConfiguration{ChatClient: chatClient}
}
