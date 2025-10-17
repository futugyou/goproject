package openai

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	rawopenai "github.com/openai/openai-go/v3"
)

// TODO: OpenAI official golang sdk has not yet implemented the responses interface
type OpenAIResponseChatClient struct {
	metadata     chatcompletion.ChatClientMetadata
	openAIClient *rawopenai.Client
	modelId      *string
}

func NewOpenAIResponseChatClient(openAIClient *rawopenai.Client, modelId *string) *OpenAIResponseChatClient {
	name := "openai"
	return &OpenAIResponseChatClient{
		metadata: chatcompletion.ChatClientMetadata{
			ProviderName: &name,
		},
		openAIClient: openAIClient,
		modelId:      modelId,
	}
}

func (client *OpenAIResponseChatClient) GetResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	return nil, nil
}

func (client *OpenAIResponseChatClient) GetStreamingResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	result := make(chan chatcompletion.ChatStreamingResponse)
	return result
}
