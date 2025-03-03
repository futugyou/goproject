package chatcompletion

import (
	"context"

	"github.com/futugyou/ai-extension/abstractions/chatcompletion"
)

type GetResponseFunc func(
	ctx context.Context,
	messages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
	client chatcompletion.IChatClient,
) (*chatcompletion.ChatResponse, error)

type GetStreamingResponseFunc func(
	ctx context.Context,
	messages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
	client chatcompletion.IChatClient,
) <-chan chatcompletion.ChatResponseUpdate

type SharedFunc func(
	messages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
	callback func(
		ctx context.Context,
		messages []chatcompletion.ChatMessage,
		options *chatcompletion.ChatOptions,
	) error,
	ctx context.Context,
) error

type AnonymousDelegatingChatClient struct {
	chatcompletion.DelegatingChatClient
	sharedFunc               SharedFunc
	getResponseFunc          GetResponseFunc
	getStreamingResponseFunc GetStreamingResponseFunc
}

func NewAnonymousDelegatingChatClient(
	innerClient chatcompletion.IChatClient,
	sharedFunc SharedFunc,
	getResponseFunc GetResponseFunc,
	getStreamingResponseFunc GetStreamingResponseFunc,
) *AnonymousDelegatingChatClient {
	return &AnonymousDelegatingChatClient{
		DelegatingChatClient: chatcompletion.DelegatingChatClient{
			InnerClient: innerClient,
		},
		sharedFunc:               sharedFunc,
		getResponseFunc:          getResponseFunc,
		getStreamingResponseFunc: getStreamingResponseFunc,
	}
}

func (client *AnonymousDelegatingChatClient) GetResponse(
	ctx context.Context,
	messages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
) (*chatcompletion.ChatResponse, error) {
	return client.getResponseFunc(ctx, messages, options, client.InnerClient)
}

func (client *AnonymousDelegatingChatClient) GetStreamingResponse(
	ctx context.Context,
	messages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
) <-chan chatcompletion.ChatResponseUpdate {
	return client.getStreamingResponseFunc(ctx, messages, options, client.InnerClient)
}
