package openai

import (
	"context"

	"github.com/futugyou/yomawari/generative-ai/abstractions/chatcompletion"
	rawopenai "github.com/openai/openai-go"
)

type OpenAIChatClient struct {
	metadata     chatcompletion.ChatClientMetadata
	openAIClient *rawopenai.Client
	modelId      *string
}

func NewOpenAIChatClient(openAIClient *rawopenai.Client, modelId *string) *OpenAIChatClient {
	name := "openai"
	return &OpenAIChatClient{
		metadata: chatcompletion.ChatClientMetadata{
			ProviderName: &name,
		},
		openAIClient: openAIClient,
		modelId:      modelId,
	}
}

func (client *OpenAIChatClient) GetResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	request := ToOpenAIChatRequest(options)
	request.Messages = rawopenai.F(ToOpenAIMessages(chatMessages))
	response, err := client.openAIClient.Chat.Completions.New(ctx, *request)
	if err != nil {
		return nil, err
	}
	return ToChatResponse(response), nil
}

func (client *OpenAIChatClient) GetStreamingResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	result := make(chan chatcompletion.ChatStreamingResponse)
	request := ToOpenAIChatRequest(options)
	request.Messages = rawopenai.F(ToOpenAIMessages(chatMessages))
	stream := client.openAIClient.Chat.Completions.NewStreaming(ctx, *request)
	go func() {
		defer close(result)
		defer stream.Close()
		for stream.Next() {
			response := stream.Current()
			ch := ToChatResponseUpdate(&response)

			select {
			case result <- chatcompletion.ChatStreamingResponse{Update: ch}:
			case <-ctx.Done():
				result <- chatcompletion.ChatStreamingResponse{Err: ctx.Err()}
				return
			}
		}

		if err := stream.Err(); err != nil {
			result <- chatcompletion.ChatStreamingResponse{Err: err}
			return
		}

	}()

	return result
}
