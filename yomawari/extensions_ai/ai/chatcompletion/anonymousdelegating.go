package chatcompletion

import (
	"context"
	"fmt"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
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
) <-chan chatcompletion.ChatStreamingResponse

type SharedFunc func(
	ctx context.Context,
	messages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
	callback func(
		ctx context.Context,
		messages []chatcompletion.ChatMessage,
		options *chatcompletion.ChatOptions,
	) error,
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
	if client.sharedFunc != nil {
		var response *chatcompletion.ChatResponse
		var err error
		client.sharedFunc(ctx, messages, options, func(
			ctx context.Context,
			messages []chatcompletion.ChatMessage,
			options *chatcompletion.ChatOptions,
		) error {
			response, err = client.InnerClient.GetResponse(ctx, messages, options)
			return err
		})

		return response, err
	}

	if client.getResponseFunc != nil {
		return client.getResponseFunc(ctx, messages, options, client.InnerClient)
	}

	if client.getStreamingResponseFunc == nil {
		return nil, fmt.Errorf("getStreamingResponseFunc is nil")
	}

	updateResponse := <-client.getStreamingResponseFunc(ctx, messages, options, client.InnerClient)
	if updateResponse.Err != nil {
		return nil, updateResponse.Err
	}

	response := chatcompletion.ToChatResponse([]chatcompletion.ChatResponseUpdate{*updateResponse.Update})
	return &response, nil
}

func (client *AnonymousDelegatingChatClient) GetStreamingResponse(
	ctx context.Context,
	messages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
) <-chan chatcompletion.ChatStreamingResponse {
	if client.getStreamingResponseFunc != nil {
		return client.getStreamingResponseFunc(ctx, messages, options, client.InnerClient)
	}

	if client.sharedFunc != nil {
		var response <-chan chatcompletion.ChatStreamingResponse
		client.sharedFunc(ctx, messages, options, func(
			ctx context.Context,
			messages []chatcompletion.ChatMessage,
			options *chatcompletion.ChatOptions,
		) error {
			response = client.InnerClient.GetStreamingResponse(ctx, messages, options)
			return nil
		})

		return response
	}

	streamResp := make(chan chatcompletion.ChatStreamingResponse)
	if client.getResponseFunc == nil {
		streamResp <- chatcompletion.ChatStreamingResponse{
			Update: nil,
			Err:    fmt.Errorf("getResponseFunc is nil"),
		}
		close(streamResp)
		return streamResp
	}

	chat, err := client.getResponseFunc(ctx, messages, options, client.InnerClient)
	if err != nil {
		streamResp <- chatcompletion.ChatStreamingResponse{
			Update: nil,
			Err:    err,
		}
		close(streamResp)
		return streamResp
	}

	updates := chat.ToChatResponseUpdates()
	go func() {
		defer close(streamResp)
		for _, item := range updates {
			streamResp <- chatcompletion.ChatStreamingResponse{
				Update: &item,
				Err:    nil,
			}
		}
	}()
	return streamResp
}
