package chatcompletion

import (
	"context"
	"fmt"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
)

type ConfigureOptionsChatClient struct {
	chatcompletion.DelegatingChatClient
	configureOptions func(*chatcompletion.ChatOptions)
}

func NewConfigureOptionsChatClient(
	innerClient chatcompletion.IChatClient,
	configureOptions func(*chatcompletion.ChatOptions),
) *ConfigureOptionsChatClient {
	return &ConfigureOptionsChatClient{
		DelegatingChatClient: chatcompletion.DelegatingChatClient{
			InnerClient: innerClient,
		},
		configureOptions: configureOptions,
	}
}

func (client *ConfigureOptionsChatClient) GetResponse(
	ctx context.Context,
	chatMessages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
) (*chatcompletion.ChatResponse, error) {
	if client.configureOptions == nil {
		return nil, fmt.Errorf("configureOptions is nil")
	}

	return client.InnerClient.GetResponse(ctx, chatMessages, client.Configure(options))
}

func (client *ConfigureOptionsChatClient) GetStreamingResponse(
	ctx context.Context,
	chatMessages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
) <-chan chatcompletion.ChatStreamingResponse {
	if client.configureOptions == nil {
		result := make(chan chatcompletion.ChatStreamingResponse, 1)
		result <- chatcompletion.ChatStreamingResponse{Err: fmt.Errorf("configureOptions is nil")}
		close(result)
		return result
	}

	return client.InnerClient.GetStreamingResponse(ctx, chatMessages, client.Configure(options))
}

func (client *ConfigureOptionsChatClient) Configure(options *chatcompletion.ChatOptions) *chatcompletion.ChatOptions {
	if options == nil {
		options = new(chatcompletion.ChatOptions)
	}

	client.configureOptions(options)

	return options
}
