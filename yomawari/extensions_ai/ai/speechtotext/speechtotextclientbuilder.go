package speechtotext

import (
	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/speechtotext"
)

type SpeechToTextClientBuilder struct {
	innerClientFactory func(core.IServiceProvider) speechtotext.ISpeechToTextClient
	clientFactories    []func(speechtotext.ISpeechToTextClient, core.IServiceProvider) speechtotext.ISpeechToTextClient
}

func NewSpeechToTextClientBuilder(client speechtotext.ISpeechToTextClient) *SpeechToTextClientBuilder {
	return &SpeechToTextClientBuilder{
		innerClientFactory: func(core.IServiceProvider) speechtotext.ISpeechToTextClient {
			return client
		},
		clientFactories: []func(speechtotext.ISpeechToTextClient, core.IServiceProvider) speechtotext.ISpeechToTextClient{},
	}
}

func NewSpeechToTextClientBuilderWithFactory(factory func(core.IServiceProvider) speechtotext.ISpeechToTextClient) *SpeechToTextClientBuilder {
	return &SpeechToTextClientBuilder{
		innerClientFactory: factory,
		clientFactories:    []func(speechtotext.ISpeechToTextClient, core.IServiceProvider) speechtotext.ISpeechToTextClient{},
	}
}

func (b *SpeechToTextClientBuilder) Build(services core.IServiceProvider) speechtotext.ISpeechToTextClient {
	chatClient := b.innerClientFactory(services)
	for i := len(b.clientFactories) - 1; i >= 0; i-- {
		factory := b.clientFactories[i]
		chatClient = factory(chatClient, services)
	}

	return chatClient
}

func (b *SpeechToTextClientBuilder) Use(factory func(speechtotext.ISpeechToTextClient, core.IServiceProvider) speechtotext.ISpeechToTextClient) *SpeechToTextClientBuilder {
	b.clientFactories = append(b.clientFactories, factory)
	return b
}

func (b *SpeechToTextClientBuilder) UseWithoutPrivider(factory func(speechtotext.ISpeechToTextClient) speechtotext.ISpeechToTextClient) *SpeechToTextClientBuilder {
	return b.Use(func(client speechtotext.ISpeechToTextClient, services core.IServiceProvider) speechtotext.ISpeechToTextClient {
		return factory(client)
	})
}

func (b *SpeechToTextClientBuilder) UseLogging(configure func(*LoggingSpeechToTextClient)) *SpeechToTextClientBuilder {
	return b.Use(func(client speechtotext.ISpeechToTextClient, services core.IServiceProvider) speechtotext.ISpeechToTextClient {
		logger, _ := core.GetService[logger.Logger](services)
		metadata, _ := core.GetService[*speechtotext.SpeechToTextClientMetadata](services)
		logclient := NewLoggingSpeechToTextClient(
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

func (b *SpeechToTextClientBuilder) ConfigureOptions(configure func(*speechtotext.SpeechToTextOptions)) *SpeechToTextClientBuilder {
	b.Use(func(innerClient speechtotext.ISpeechToTextClient, sp core.IServiceProvider) speechtotext.ISpeechToTextClient {
		return NewConfigureOptionsSpeechToTextClient(innerClient, configure, nil)
	})
	return b
}

func ISpeechToTextClientAsBuilder(innerClient speechtotext.ISpeechToTextClient) *SpeechToTextClientBuilder {
	return NewSpeechToTextClientBuilder(innerClient)
}
