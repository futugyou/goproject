package openai

import (
	"context"
	"encoding/json"

	"github.com/futugyou/yomawari/extensions_ai/abstractions"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
	rawopenai "github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/openai/openai-go/shared"
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

	for msg := range originalResponse {
		if msg.Err != nil {
			return nil, msg.Err
		}
		updates = append(updates, *msg.Update)
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
		params := getBetaThreadNewAndRunParams(chatMessages, options)
		params.AssistantID = rawopenai.F(client.assistantId)
		if len(params.Model.Value) == 0 {
			params.Model = rawopenai.F(*options.ModelId)
		}

		stream = client.threads.NewAndRunStreaming(ctx, params)
	} else {
		params := getBetaThreadRunNewParams(chatMessages, options)
		params.AssistantID = rawopenai.F(client.assistantId)
		if len(params.Model.Value) == 0 {
			params.Model = rawopenai.F(*options.ModelId)
		}
		stream = client.threads.Runs.NewStreaming(ctx, threadId, params)
	}

	var newThreadId *string = &threadId
	var responseId *string
	var modelId *string = options.ModelId
	go func() {
		defer close(result)
		defer stream.Close()
		for stream.Next() {
			response := stream.Current()
			update, err := ToChatResponseUpdateFromAssistantStreamEvent(response, &newThreadId, &responseId, &modelId)
			if err != nil {
				result <- chatcompletion.ChatStreamingResponse{Err: err}
				return
			}

			if update == nil || update.ResponseId == nil {
				continue
			}

			select {
			case result <- chatcompletion.ChatStreamingResponse{
				Update: update,
			}:
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

func getBetaThreadRunNewParams(chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) rawopenai.BetaThreadRunNewParams {
	result := rawopenai.BetaThreadRunNewParams{
		Include: rawopenai.F([]rawopenai.RunStepInclude{rawopenai.RunStepIncludeStepDetailsToolCallsFileSearchResultsContent}),
	}

	if options != nil {
		if options.TopP != nil && *options.TopP > 0 {
			result.TopP = rawopenai.F(*options.TopP)
		}
		if options.MaxOutputTokens != nil && *options.MaxOutputTokens > 0 {
			result.MaxCompletionTokens = rawopenai.F(*options.MaxOutputTokens)
		}
		if options.Temperature != nil && *options.Temperature > 0 {
			result.Temperature = rawopenai.F(*options.Temperature)
		}
		if options.ModelId != nil && len(*options.ModelId) > 0 {
			result.Model = rawopenai.F(*options.ModelId)
		}
		if options.ToolMode != nil && len(*options.ToolMode) > 0 {
			var toolChoice rawopenai.AssistantToolChoiceOptionUnionParam
			switch *options.ToolMode {
			case "requireAny":
				toolChoice = rawopenai.AssistantToolChoiceOptionAutoRequired
			case "auto":
				toolChoice = rawopenai.AssistantToolChoiceOptionAutoAuto
			case "none":
				toolChoice = rawopenai.AssistantToolChoiceOptionAutoNone
			}
			if toolChoice != nil {
				result.ToolChoice = rawopenai.F(toolChoice)
			}
		}
		if v, ok := options.AdditionalProperties["ParallelToolCalls"].(bool); ok {
			result.ParallelToolCalls = rawopenai.F(v)
		}
		if v, ok := options.AdditionalProperties["MaxPromptTokens"].(int64); ok {
			result.MaxPromptTokens = rawopenai.F(v)
		}
		if v, ok := options.AdditionalProperties["TruncationStrategy"].(rawopenai.BetaThreadRunNewParamsTruncationStrategy); ok {
			result.TruncationStrategy = rawopenai.F(v)
		}
		tools := []rawopenai.AssistantToolUnionParam{}
		for _, tool := range options.Tools {
			if t, ok := tool.(functions.AIFunction); ok {
				var m shared.FunctionParameters = t.GetJsonSchema()
				strict := false
				if v, ok := t.GetAdditionalProperties()["Strict"].(bool); ok {
					strict = v
				}
				pa := shared.FunctionDefinitionParam{
					Name:        rawopenai.F(t.GetName()),
					Description: rawopenai.F(t.GetDescription()),
					Strict:      rawopenai.F(strict),
					Parameters:  rawopenai.F(m),
				}

				tools = append(tools, rawopenai.FunctionToolParam{
					Function: rawopenai.F(pa),
					Type:     rawopenai.F(rawopenai.FunctionToolTypeFunction),
				})
				continue
			}
			if _, ok := tool.(abstractions.CodeInterpreterTool); ok {
				tools = append(tools, rawopenai.CodeInterpreterToolParam{
					Type: rawopenai.F(rawopenai.CodeInterpreterToolTypeCodeInterpreter),
				})
			}
		}
		if len(tools) > 0 {
			result.Tools = rawopenai.F(tools)
		}
	}

	instructions := ""
	message := []rawopenai.BetaThreadRunNewParamsAdditionalMessage{}
	for _, msg := range chatMessages {
		if msg.Role == "system" || msg.Role == "developer" {
			for _, con := range msg.Contents {
				if con, ok := con.(*contents.TextContent); ok {
					instructions += con.Text
				}
			}
			continue
		}

		for _, con := range msg.Contents {
			role := rawopenai.BetaThreadRunNewParamsAdditionalMessagesRoleUser
			if msg.Role == "assistant" {
				role = rawopenai.BetaThreadRunNewParamsAdditionalMessagesRoleAssistant
			}
			if con, ok := con.(*contents.TextContent); ok {
				var sss = []rawopenai.MessageContentPartParamUnion{
					rawopenai.MessageContentPartParam{
						Type: rawopenai.F(rawopenai.MessageContentPartParamTypeText),
						Text: rawopenai.F(con.Text),
					},
				}
				var content []rawopenai.MessageContentPartParamUnion = sss

				message = append(message, rawopenai.BetaThreadRunNewParamsAdditionalMessage{
					Role:    rawopenai.F(role),
					Content: rawopenai.F(content),
				})
				continue
			}
			if con, ok := con.(*contents.DataContent); ok && con.MediaTypeStartsWith("image") {
				var sss = []rawopenai.MessageContentPartParamUnion{
					rawopenai.MessageContentPartParam{
						Type: rawopenai.F(rawopenai.MessageContentPartParamTypeImageURL),
						ImageURL: rawopenai.F(rawopenai.ImageURLParam{
							URL: rawopenai.F(con.GetURI()),
						}),
					},
				}
				var content []rawopenai.MessageContentPartParamUnion = sss

				message = append(message, rawopenai.BetaThreadRunNewParamsAdditionalMessage{
					Role:    rawopenai.F(role),
					Content: rawopenai.F(content),
				})
			}
		}
	}

	if len(instructions) > 0 {
		result.AdditionalInstructions = rawopenai.F(instructions)
	}

	if len(message) > 0 {
		result.AdditionalMessages = rawopenai.F(message)
	}

	return result
}

func getBetaThreadNewAndRunParams(chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) rawopenai.BetaThreadNewAndRunParams {
	result := rawopenai.BetaThreadNewAndRunParams{}

	if options != nil {
		if options.TopP != nil && *options.TopP > 0 {
			result.TopP = rawopenai.F(*options.TopP)
		}
		if options.MaxOutputTokens != nil && *options.MaxOutputTokens > 0 {
			result.MaxCompletionTokens = rawopenai.F(*options.MaxOutputTokens)
		}
		if options.Temperature != nil && *options.Temperature > 0 {
			result.Temperature = rawopenai.F(*options.Temperature)
		}
		if options.ModelId != nil && len(*options.ModelId) > 0 {
			result.Model = rawopenai.F(*options.ModelId)
		}
		if options.ToolMode != nil && len(*options.ToolMode) > 0 {
			var toolChoice rawopenai.AssistantToolChoiceOptionUnionParam
			switch *options.ToolMode {
			case "requireAny":
				toolChoice = rawopenai.AssistantToolChoiceOptionAutoRequired
			case "auto":
				toolChoice = rawopenai.AssistantToolChoiceOptionAutoAuto
			case "none":
				toolChoice = rawopenai.AssistantToolChoiceOptionAutoNone
			}
			if toolChoice != nil {
				result.ToolChoice = rawopenai.F(toolChoice)
			}
		}
		if v, ok := options.AdditionalProperties["ParallelToolCalls"].(bool); ok {
			result.ParallelToolCalls = rawopenai.F(v)
		}
		if v, ok := options.AdditionalProperties["MaxPromptTokens"].(int64); ok {
			result.MaxPromptTokens = rawopenai.F(v)
		}
		if v, ok := options.AdditionalProperties["TruncationStrategy"].(rawopenai.BetaThreadNewAndRunParamsTruncationStrategy); ok {
			result.TruncationStrategy = rawopenai.F(v)
		}
		tools := []rawopenai.BetaThreadNewAndRunParamsToolUnion{}
		for _, tool := range options.Tools {
			if t, ok := tool.(functions.AIFunction); ok {
				var m shared.FunctionParameters = t.GetJsonSchema()
				strict := false
				if v, ok := t.GetAdditionalProperties()["strictJsonSchema"].(bool); ok {
					strict = v
				}
				pa := shared.FunctionDefinitionParam{
					Name:        rawopenai.F(t.GetName()),
					Description: rawopenai.F(t.GetDescription()),
					Strict:      rawopenai.F(strict),
					Parameters:  rawopenai.F(m),
				}

				tools = append(tools, rawopenai.FunctionToolParam{
					Function: rawopenai.F(pa),
					Type:     rawopenai.F(rawopenai.FunctionToolTypeFunction),
				})
				continue
			}
			if _, ok := tool.(abstractions.CodeInterpreterTool); ok {
				tools = append(tools, rawopenai.CodeInterpreterToolParam{
					Type: rawopenai.F(rawopenai.CodeInterpreterToolTypeCodeInterpreter),
				})
			}
		}
		if len(tools) > 0 {
			result.Tools = rawopenai.F(tools)
		}
	}

	message := []rawopenai.BetaThreadNewAndRunParamsThreadMessage{}
	for _, msg := range chatMessages {
		if msg.Role == "system" || msg.Role == "developer" {
			continue
		}
		for _, con := range msg.Contents {
			role := rawopenai.BetaThreadNewAndRunParamsThreadMessagesRoleUser
			if msg.Role == "assistant" {
				role = rawopenai.BetaThreadNewAndRunParamsThreadMessagesRoleAssistant
			}
			if con, ok := con.(*contents.TextContent); ok {
				var sss = []rawopenai.MessageContentPartParamUnion{
					rawopenai.MessageContentPartParam{
						Type: rawopenai.F(rawopenai.MessageContentPartParamTypeText),
						Text: rawopenai.F(con.Text),
					},
				}

				var content []rawopenai.MessageContentPartParamUnion = sss
				message = append(message, rawopenai.BetaThreadNewAndRunParamsThreadMessage{
					Role:    rawopenai.F(role),
					Content: rawopenai.F(content),
				})

				continue
			}
			if con, ok := con.(*contents.DataContent); ok && con.MediaTypeStartsWith("image") {
				var sss = []rawopenai.MessageContentPartParamUnion{
					rawopenai.MessageContentPartParam{
						Type: rawopenai.F(rawopenai.MessageContentPartParamTypeImageURL),
						ImageURL: rawopenai.F(rawopenai.ImageURLParam{
							URL: rawopenai.F(con.GetURI()),
						}),
					},
				}

				var content []rawopenai.MessageContentPartParamUnion = sss
				message = append(message, rawopenai.BetaThreadNewAndRunParamsThreadMessage{
					Role:    rawopenai.F(role),
					Content: rawopenai.F(content),
				})
			}
		}
	}

	if len(message) > 0 {
		result.Thread = rawopenai.F(rawopenai.BetaThreadNewAndRunParamsThread{
			Messages: rawopenai.F(message),
		})
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
