package openai

import (
	"context"
	"encoding/json"

	"github.com/futugyou/yomawari/generative-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/generative-ai/abstractions/contents"
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
	originalResponse := client.GetStreamingResponse(ctx, chatMessages, options)
	updates := []chatcompletion.ChatResponseUpdate{}
	var err error
	go func() {
		for msg := range originalResponse {
			if msg.Err != nil {
				err = msg.Err
				return
			}
			updates = append(updates, *msg.Update)
		}
	}()
	if err != nil {
		return nil, err
	}
	chatResponse := chatcompletion.ToChatResponse(updates)
	return &chatResponse, nil
}

func (client *OpenAIAssistantClient) GetStreamingResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	result := make(chan chatcompletion.ChatStreamingResponse)
	threadId := getThreadId(client.threadId, options)
	runId, tools := getFunctionResultContents(chatMessages)

	var stream *ssestream.Stream[rawopenai.AssistantStreamEvent]
	if runId != nil && len(*runId) > 0 && len(threadId) > 0 {
		params := rawopenai.BetaThreadRunSubmitToolOutputsParams{
			ToolOutputs: rawopenai.F(tools),
		}
		stream = client.threads.Runs.SubmitToolOutputsStreaming(ctx, threadId, *runId, params)
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

func getThreadId(threadId *string, options *chatcompletion.ChatOptions) string {
	if options != nil && options.ChatThreadId != nil && len(*options.ChatThreadId) > 0 {
		return *options.ChatThreadId
	}

	if threadId != nil && len(*threadId) > 0 {
		return *threadId
	}

	return ""
}

func getFunctionResultContents(messages []chatcompletion.ChatMessage) (*string, []rawopenai.BetaThreadRunSubmitToolOutputsParamsToolOutput) {
	var existingRunId *string
	tools := []rawopenai.BetaThreadRunSubmitToolOutputsParamsToolOutput{}

	for _, message := range messages {
		for _, con := range message.Contents {
			if c, ok := con.(*contents.FunctionResultContent); ok {
				var strSlice []string

				if err := json.Unmarshal([]byte(c.CallId), &strSlice); err != nil {
					continue
				}

				if len(strSlice) != 2 || len(strSlice[0]) == 0 || len(strSlice[1]) == 0 || (existingRunId != nil && *existingRunId != strSlice[0]) {
					continue
				}
				existingRunId = &strSlice[0]
				result := ""
				if r, ok := c.Result.(string); ok {
					result = r
				}
				tools = append(tools, rawopenai.BetaThreadRunSubmitToolOutputsParamsToolOutput{
					ToolCallID: rawopenai.F(strSlice[1]),
					Output:     rawopenai.F(result),
				})
			}
		}
	}

	return existingRunId, tools
}
