package chatcompletion

import (
	"github.com/futugyou/ai-extension/abstractions/chatcompletion"
	"github.com/futugyou/ai-extension/core"
)

// TODO, find a DI container
// [dig](https://github.com/uber-go/dig) or [wire](https://github.com/google/wire)
type ChatClientBuilder struct {
	innerClientFactory func(core.IServiceProvider) chatcompletion.IChatClient
	clientFactories    []func(chatcompletion.IChatClient, core.IServiceProvider) chatcompletion.IChatClient
}

func NewChatClientBuilder(client chatcompletion.IChatClient) *ChatClientBuilder {
	return &ChatClientBuilder{
		innerClientFactory: func(core.IServiceProvider) chatcompletion.IChatClient {
			return client
		},
		clientFactories: []func(chatcompletion.IChatClient, core.IServiceProvider) chatcompletion.IChatClient{},
	}
}

func NewChatClientBuilderWithFactory(factory func(core.IServiceProvider) chatcompletion.IChatClient) *ChatClientBuilder {
	return &ChatClientBuilder{
		innerClientFactory: factory,
		clientFactories:    []func(chatcompletion.IChatClient, core.IServiceProvider) chatcompletion.IChatClient{},
	}
}

func (b *ChatClientBuilder) Build(services core.IServiceProvider) chatcompletion.IChatClient {
	chatClient := b.innerClientFactory(services)
	for i := len(b.clientFactories) - 1; i >= 0; i-- {
		factory := b.clientFactories[i]
		chatClient = factory(chatClient, services)
	}

	return chatClient
}

func (b *ChatClientBuilder) Use(factory func(chatcompletion.IChatClient, core.IServiceProvider) chatcompletion.IChatClient) *ChatClientBuilder {
	b.clientFactories = append(b.clientFactories, factory)
	return b
}

func (b *ChatClientBuilder) UseWithoutPrivider(factory func(chatcompletion.IChatClient) chatcompletion.IChatClient) *ChatClientBuilder {
	return b.Use(func(client chatcompletion.IChatClient, services core.IServiceProvider) chatcompletion.IChatClient {
		return factory(client)
	})
}

func (b *ChatClientBuilder) UseSharedFunc(sharedFunc SharedFunc) *ChatClientBuilder {
	return b.Use(func(client chatcompletion.IChatClient, services core.IServiceProvider) chatcompletion.IChatClient {
		return NewAnonymousDelegatingChatClient(
			client,
			sharedFunc,
			nil,
			nil,
		)
	})
}

func (b *ChatClientBuilder) UseResponseFunc(getResponseFunc GetResponseFunc, getStreamingResponseFunc GetStreamingResponseFunc) *ChatClientBuilder {
	return b.Use(func(client chatcompletion.IChatClient, services core.IServiceProvider) chatcompletion.IChatClient {
		return NewAnonymousDelegatingChatClient(
			client,
			nil,
			getResponseFunc,
			getStreamingResponseFunc,
		)
	})
}
