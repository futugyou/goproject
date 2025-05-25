package chatcompletion

import (
	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
)

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

func (b *ChatClientBuilder) ConfigureOptions(configure func(*chatcompletion.ChatOptions)) *ChatClientBuilder {
	b.Use(func(innerClient chatcompletion.IChatClient, sp core.IServiceProvider) chatcompletion.IChatClient {
		return NewConfigureOptionsChatClient(innerClient, configure)
	})
	return b
}

func (b *ChatClientBuilder) UseDistributedCache(storage core.IDistributedCache, configure func(*DistributedCachingChatClient)) *ChatClientBuilder {
	b.Use(func(innerClient chatcompletion.IChatClient, sp core.IServiceProvider) chatcompletion.IChatClient {
		if storage == nil {
			storage, _ = core.GetService[core.IDistributedCache](sp)
		}

		var chatClient = NewDistributedCachingChatClient(innerClient, storage)
		if configure != nil {
			configure(chatClient)
		}

		return chatClient
	})
	return b
}

func (b *ChatClientBuilder) UseFunctionInvocation(configure func(*FunctionInvokingChatClient)) *ChatClientBuilder {
	b.Use(func(innerClient chatcompletion.IChatClient, sp core.IServiceProvider) chatcompletion.IChatClient {
		var chatClient = NewFunctionInvokingChatClient(innerClient)
		if configure != nil {
			configure(chatClient)
		}

		return chatClient
	})
	return b
}

func (b *ChatClientBuilder) UseLogging(configure func(*LoggingChatClient)) *ChatClientBuilder {
	return b.Use(func(client chatcompletion.IChatClient, services core.IServiceProvider) chatcompletion.IChatClient {
		logger, _ := core.GetService[logger.Logger](services)
		metadata, _ := core.GetService[*chatcompletion.ChatClientMetadata](services)
		logclient := NewLoggingChatClient(
			client,
			logger,
			metadata,
		)

		if configure != nil {
			configure(logclient)
		}

		return logclient
	})
}

func (b *ChatClientBuilder) UseOpenTelemetry(configure func(*OpenTelemetryChatClient)) *ChatClientBuilder {
	return b.Use(func(client chatcompletion.IChatClient, services core.IServiceProvider) chatcompletion.IChatClient {
		metadata, _ := core.GetService[*chatcompletion.ChatClientMetadata](services)
		otelclient := NewOpenTelemetryChatClient(
			client,
			metadata,
		)

		if configure != nil {
			configure(otelclient)
		}

		return otelclient
	})
}

func IChatClientAsBuilder(innerClient chatcompletion.IChatClient) *ChatClientBuilder {
	return NewChatClientBuilder(innerClient)
}
