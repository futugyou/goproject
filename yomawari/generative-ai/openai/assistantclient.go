package openai

import (
	"context"

	"github.com/futugyou/yomawari/generative-ai/abstractions/chatcompletion"
	rawopenai "github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
)

type OpenAIAssistantClient struct {
	metadata    chatcompletion.ChatClientMetadata
	threads     *rawopenai.BetaThreadService
	assistantId string
	threadId    *string
}

func NewOpenAIAssistantClient(threads *rawopenai.BetaThreadService, assistantId string, threadId *string) *OpenAIAssistantClient {
	name := "openai"
	return &OpenAIAssistantClient{
		metadata: chatcompletion.ChatClientMetadata{
			ProviderName: &name,
		},
		threads:     threads,
		assistantId: assistantId,
		threadId:    threadId,
	}
}

func (client *OpenAIAssistantClient) GetResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {

	return nil, nil
}

func (client *OpenAIAssistantClient) GetStreamingResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	result := make(chan chatcompletion.ChatStreamingResponse)
	threadId := ""
	runId := ""
	//runId = getRunId(....)
	if options != nil && options.ChatThreadId != nil {
		threadId = *options.ChatThreadId
	}

	var stream *ssestream.Stream[rawopenai.AssistantStreamEvent]
	if len(runId) > 0 && len(threadId) > 0 {
		// TODO
		params := rawopenai.BetaThreadRunSubmitToolOutputsParams{
			ToolOutputs: rawopenai.F([]rawopenai.BetaThreadRunSubmitToolOutputsParamsToolOutput{}),
		}
		stream = client.threads.Runs.SubmitToolOutputsStreaming(ctx, threadId, runId, params)
	} else if len(threadId) == 0 {
		// TODO
		params := rawopenai.BetaThreadNewAndRunParams{}
		stream = client.threads.NewAndRunStreaming(ctx, params)
	} else {
		// TODO
		params := rawopenai.BetaThreadRunNewParams{}
		stream = client.threads.Runs.NewStreaming(ctx, threadId, params)
	}
	for stream.Next() {
		evt := stream.Current()
		result <- chatcompletion.ChatStreamingResponse{
			Update: ToChatResponseUpdateFromAssistantStreamEvent(evt),
		}
	}
	return result
}
