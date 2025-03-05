package chatcompletion

import (
	"github.com/futugyou/ai-extension/abstractions/chatcompletion"
)

// TODO, find a DI container
// [dig](https://github.com/uber-go/dig) or [wire](https://github.com/google/wire)
type ChatClientBuilder struct {
	innerClientFactory func(interface{}) chatcompletion.IChatClient
	clientFactories    []func(chatcompletion.IChatClient, interface{}) chatcompletion.IChatClient
}

func NewChatClientBuilder(client chatcompletion.IChatClient) *ChatClientBuilder {
	return &ChatClientBuilder{
		innerClientFactory: func(interface{}) chatcompletion.IChatClient {
			return client
		},
		clientFactories: []func(chatcompletion.IChatClient, interface{}) chatcompletion.IChatClient{},
	}
}

func NewChatClientBuilderWithFactory(factory func(interface{}) chatcompletion.IChatClient) *ChatClientBuilder {
	return &ChatClientBuilder{
		innerClientFactory: factory,
		clientFactories:    []func(chatcompletion.IChatClient, interface{}) chatcompletion.IChatClient{},
	}
}

func (b *ChatClientBuilder) Build(services interface{}) chatcompletion.IChatClient {
	chatClient := b.innerClientFactory(services)
	for i := len(b.clientFactories) - 1; i >= 0; i-- {
		factory := b.clientFactories[i]
		chatClient = factory(chatClient, services)
	}

	return chatClient
}

func (b *ChatClientBuilder) Use(factory func(chatcompletion.IChatClient, interface{}) chatcompletion.IChatClient) *ChatClientBuilder {
	b.clientFactories = append(b.clientFactories, factory)
	return b
}

func (b *ChatClientBuilder) UseWithoutPrivider(factory func(chatcompletion.IChatClient) chatcompletion.IChatClient) *ChatClientBuilder {
	return b.Use(func(client chatcompletion.IChatClient, services interface{}) chatcompletion.IChatClient {
		return factory(client)
	})
}

func (b *ChatClientBuilder) UseSharedFunc(sharedFunc SharedFunc) *ChatClientBuilder {
	return b.Use(func(client chatcompletion.IChatClient, services interface{}) chatcompletion.IChatClient {
		return NewAnonymousDelegatingChatClient(
			client,
			sharedFunc,
			nil,
			nil,
		)
	})
}

func (b *ChatClientBuilder) UseResponseFunc(getResponseFunc GetResponseFunc, getStreamingResponseFunc GetStreamingResponseFunc) *ChatClientBuilder {
	return b.Use(func(client chatcompletion.IChatClient, services interface{}) chatcompletion.IChatClient {
		return NewAnonymousDelegatingChatClient(
			client,
			nil,
			getResponseFunc,
			getStreamingResponseFunc,
		)
	})
}
