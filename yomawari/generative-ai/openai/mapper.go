package openai

import (
	"encoding/json"
	"time"

	"github.com/futugyou/yomawari/generative-ai/abstractions"
	"github.com/futugyou/yomawari/generative-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/generative-ai/abstractions/contents"
	"github.com/futugyou/yomawari/generative-ai/abstractions/embeddings"
	"github.com/futugyou/yomawari/generative-ai/abstractions/functions"
	rawopenai "github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

// Conversion functions
func ToOpenAIChatCompletion(response chatcompletion.ChatResponse, options json.RawMessage) rawopenai.ChatCompletion {
	result := rawopenai.ChatCompletion{}
	return result
}

func ToOpenAIMessages(chatMessages []chatcompletion.ChatMessage) []rawopenai.ChatCompletionMessageParamUnion {
	result := []rawopenai.ChatCompletionMessageParamUnion{}
	for _, msg := range chatMessages {
		result = append(result, ToOpenAIMessage(msg))
	}

	return result
}

func ToOpenAIMessage(chatMessage chatcompletion.ChatMessage) rawopenai.ChatCompletionMessageParamUnion {
	result := rawopenai.ChatCompletionMessageParam{}

	var role string = (string)(chatMessage.Role)
	result.Role = rawopenai.F(rawopenai.ChatCompletionMessageParamRole(role))
	if chatMessage.AuthorName != nil {
		result.Name = rawopenai.F(*chatMessage.AuthorName)
	}
	result.Content = rawopenai.F(ToOpenAIContents(chatMessage.Contents, role))
	if v, ok := chatMessage.AdditionalProperties["Refusal"].(string); ok {
		result.Refusal = rawopenai.F(v)
	}

	if role == "assistant" {
		tools := []rawopenai.ChatCompletionMessageToolCallParam{}
		for _, con := range chatMessage.Contents {
			if c, ok := con.(*contents.FunctionCallContent); ok {
				if d, err := json.Marshal(c.Arguments); err != nil {
					sf := rawopenai.ChatCompletionMessageToolCallParam{
						ID: rawopenai.F(c.CallId),
						Function: rawopenai.F(rawopenai.ChatCompletionMessageToolCallFunctionParam{
							Arguments: rawopenai.F(string(d)),
							Name:      rawopenai.F(c.Name),
						}),
						Type: rawopenai.F(rawopenai.ChatCompletionMessageToolCallTypeFunction),
					}
					tools = append(tools, sf)
				}
			}
		}
		var tool interface{} = tools
		result.ToolCalls = rawopenai.F(tool)
	}

	return result
}

func ToOpenAIContents(cons []contents.IAIContent, role string) interface{} {
	result := []rawopenai.ChatCompletionContentPartUnionParam{}
	for _, v := range cons {
		switch con := v.(type) {
		case *contents.TextContent:
			result = append(result, rawopenai.TextPart(con.Text))
		case *contents.UriContent:
			if con.MediaTypeStartsWith("image") {
				result = append(result, rawopenai.ImagePart(con.URI))
			}
		case *contents.DataContent:
			if con.MediaTypeStartsWith("image") {
				result = append(result, rawopenai.ImagePart(con.GetURI()))
			}
			if con.MediaTypeStartsWith("audio") {
				var format rawopenai.ChatCompletionContentPartInputAudioInputAudioFormat
				if con.MediaTypeStartsWith("audio/mpeg") {
					format = rawopenai.ChatCompletionContentPartInputAudioInputAudioFormatMP3
				} else if con.MediaTypeStartsWith("audio/wav") {
					format = rawopenai.ChatCompletionContentPartInputAudioInputAudioFormatWAV
				}

				if len(format) == 0 {
					break
				}

				audio := rawopenai.ChatCompletionContentPartInputAudioParam{
					InputAudio: rawopenai.F(rawopenai.ChatCompletionContentPartInputAudioInputAudioParam{
						Data:   rawopenai.F(con.GetURI()),
						Format: rawopenai.F(format),
					}),
					Type: rawopenai.F(rawopenai.ChatCompletionContentPartInputAudioTypeInputAudio),
				}
				result = append(result, audio)
			}
		}
	}
	return result
}

func ToOpenAIChatRequest(options *chatcompletion.ChatOptions) *rawopenai.ChatCompletionNewParams {
	result := &rawopenai.ChatCompletionNewParams{}
	if options == nil {
		return result
	}
	if options.ModelId != nil {
		result.Model = rawopenai.F(*options.ModelId)
	}
	if options.FrequencyPenalty != nil {
		result.FrequencyPenalty = rawopenai.F(*options.FrequencyPenalty)
	}
	if options.MaxOutputTokens != nil {
		result.MaxCompletionTokens = rawopenai.F(*options.MaxOutputTokens)
	}
	if options.PresencePenalty != nil {
		result.PresencePenalty = rawopenai.F(*options.PresencePenalty)
	}
	if options.Seed != nil {
		result.Seed = rawopenai.F(*options.Seed)
	}
	if options.Temperature != nil {
		result.Temperature = rawopenai.F(*options.Temperature)
	}
	if options.TopP != nil {
		result.TopP = rawopenai.F(*options.TopP)
	}

	var defaultFormat rawopenai.ChatCompletionNewParamsResponseFormatUnion = rawopenai.ChatCompletionNewParamsResponseFormat{
		Type: rawopenai.F(rawopenai.ChatCompletionNewParamsResponseFormatTypeJSONObject),
	}
	// TODO: json schema is not implement
	if options.ResponseFormat != nil && *options.ResponseFormat == chatcompletion.TextFormat {
		defaultFormat = rawopenai.ChatCompletionNewParamsResponseFormatUnion(rawopenai.ChatCompletionNewParamsResponseFormat{
			Type: rawopenai.F(rawopenai.ChatCompletionNewParamsResponseFormatTypeText),
		})
	}

	result.ResponseFormat = rawopenai.F(defaultFormat)

	if len(options.StopSequences) > 0 {
		var stop rawopenai.ChatCompletionNewParamsStopUnion = rawopenai.ChatCompletionNewParamsStopArray(options.StopSequences)
		result.Stop = rawopenai.F(stop)
	}

	for _, tool := range options.Tools {
		tools := []rawopenai.ChatCompletionToolParam{}
		if v, ok := tool.(functions.AIFunction); ok {
			t := ToOpenAIChatCompletionToolParam(v)
			tools = append(tools, t)
		}
		result.Tools = rawopenai.F(tools)
	}

	if options.ToolMode != nil {
		var choice rawopenai.ChatCompletionToolChoiceOptionUnionParam = rawopenai.ChatCompletionToolChoiceOptionAutoNone
		switch *options.ToolMode {
		case chatcompletion.AutoMode:
			choice = rawopenai.ChatCompletionToolChoiceOptionAutoAuto
		case chatcompletion.RequireAnyMode:
			choice = rawopenai.ChatCompletionToolChoiceOptionAutoRequired
		}
		result.ToolChoice = rawopenai.F(choice)
	}

	if v, ok := options.AdditionalProperties["ParallelToolCalls"].(bool); ok {
		result.ParallelToolCalls = rawopenai.F(v)
	}
	if v, ok := options.AdditionalProperties["Audio"].(rawopenai.ChatCompletionAudioParam); ok {
		result.Audio = rawopenai.F(v)
	}
	if v, ok := options.AdditionalProperties["User"].(string); ok {
		result.User = rawopenai.F(v)
	}
	if v, ok := options.AdditionalProperties["LogitBias"].(map[string]int64); ok {
		result.LogitBias = rawopenai.F(v)
	}
	if v, ok := options.AdditionalProperties["Metadata"].(shared.MetadataParam); ok {
		result.Metadata = rawopenai.F(v)
	}
	if v, ok := options.AdditionalProperties["Prediction"].(rawopenai.ChatCompletionPredictionContentParam); ok {
		result.Prediction = rawopenai.F(v)
	}
	if v, ok := options.AdditionalProperties["ReasoningEffort"].(rawopenai.ChatCompletionReasoningEffort); ok {
		result.ReasoningEffort = rawopenai.F(v)
	}
	if v, ok := options.AdditionalProperties["Modalities"].([]rawopenai.ChatCompletionModality); ok {
		result.Modalities = rawopenai.F(v)
	}
	if v, ok := options.AdditionalProperties["Store"].(bool); ok {
		result.Store = rawopenai.F(v)
	}
	if v, ok := options.AdditionalProperties["TopLogprobs"].(int64); ok {
		result.TopLogprobs = rawopenai.F(v)
	}
	return result
}

func ToOpenAIChatCompletionToolParam(v functions.AIFunction) rawopenai.ChatCompletionToolParam {
	pa := shared.FunctionDefinitionParam{
		Name:        rawopenai.F(v.GetName()),
		Description: rawopenai.F(v.GetDescription()),
		Strict:      rawopenai.F(false),
	}
	var m shared.FunctionParameters = v.GetAdditionalProperties()
	pa.Parameters = rawopenai.F(m)

	p := rawopenai.ChatCompletionToolParam{
		Function: rawopenai.F(pa),
		Type:     rawopenai.F(rawopenai.ChatCompletionToolTypeFunction),
	}

	return p
}

func ToChatResponse(openAICompletion *rawopenai.ChatCompletion) *chatcompletion.ChatResponse {
	if openAICompletion == nil {
		return nil
	}

	chatMessages := []chatcompletion.ChatMessage{}
	for _, v := range openAICompletion.Choices {
		chatMessages = append(chatMessages, ToChatMessage(v.Message))
	}

	chatResponse := chatcompletion.NewChatResponse(chatMessages, nil)
	return chatResponse
}

func ToChatMessage(v rawopenai.ChatCompletionMessage) chatcompletion.ChatMessage {
	message := chatcompletion.ChatMessage{
		Role:                 ToChatRole(v.Role),
		Contents:             []contents.IAIContent{},
		RawRepresentation:    v,
		AdditionalProperties: map[string]interface{}{},
	}

	if len(v.Content) > 0 {
		con := contents.TextContent{
			AIContent: contents.AIContent{
				AdditionalProperties: map[string]interface{}{},
			},
			Text: v.Content,
		}
		message.Contents = append(message.Contents, con)
	}

	if len(v.Audio.Data) > 0 {
		con := contents.DataContent{
			AIContent: contents.AIContent{
				AdditionalProperties: map[string]interface{}{},
			},
			URI:       "",
			MediaType: "audio/mpeg",
			Data:      []byte(v.Audio.Data),
		}
		con.AdditionalProperties["Id"] = v.Audio.ID
		con.AdditionalProperties["ExpiresAt"] = v.Audio.ExpiresAt
		con.AdditionalProperties["Transcript"] = v.Audio.Transcript
		message.Contents = append(message.Contents, con)
	}

	for _, tool := range v.ToolCalls {
		if len(tool.Function.Arguments) > 0 {
			con := contents.CreateFromParsedArguments(tool.Function.Arguments, tool.ID, tool.Function.Name, func(args string) (map[string]interface{}, error) {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(args), &result)
				return result, err
			})
			con.RawRepresentation = tool
			message.Contents = append(message.Contents, con)
		}
	}

	return message
}

func ToChatRole(v rawopenai.ChatCompletionMessageRole) chatcompletion.ChatRole {
	role := (string)(v)
	switch role {
	case "system":
		return chatcompletion.RoleSystem
	case "assistant":
		return chatcompletion.RoleAssistant
	case "user":
		return chatcompletion.RoleUser
	case "tool":
		return chatcompletion.RoleTool
	default:
		return chatcompletion.RoleSystem
	}
}

func ToChatResponseUpdate(response *rawopenai.ChatCompletionChunk) *chatcompletion.ChatResponseUpdate {
	if response == nil {
		return nil
	}

	created := time.Unix(response.Created, 0)
	result := &chatcompletion.ChatResponseUpdate{
		ResponseId:           &response.ID,
		ModelId:              &response.Model,
		RawRepresentation:    response,
		AdditionalProperties: map[string]interface{}{},
		Contents:             []contents.IAIContent{},
		CreatedAt:            &created,
	}

	if len(response.SystemFingerprint) > 0 {
		result.AdditionalProperties["SystemFingerprint"] = response.SystemFingerprint
	}

	if response.Usage.CompletionTokens > 0 {
		result.Contents = append(result.Contents, contents.UsageContent{
			AIContent: contents.AIContent{},
			Details: abstractions.UsageDetails{
				InputTokenCount:      &response.Usage.PromptTokens,
				OutputTokenCount:     &response.Usage.CompletionTokens,
				TotalTokenCount:      &response.Usage.TotalTokens,
				AdditionalProperties: map[string]int64{},
			},
		})
	}

	if len(response.Choices) == 0 {
		return result
	}

	finishReason := chatcompletion.ChatFinishReason((string)(response.Choices[len(response.Choices)-1].FinishReason))
	role := chatcompletion.StringToChatRole((string)(response.Choices[len(response.Choices)-1].Delta.Role))

	result.Role = &role
	result.FinishReason = &finishReason

	for _, chunk := range response.Choices {
		result.Contents = append(result.Contents, ToAIContent(chunk)...)
	}
	return result
}

type InnerContentStruct struct {
	Type       string                  `json:"type"`
	Image      InnerContentImageStruct `json:"image_url"`
	Refusal    string                  `json:"refusal"`
	InputAudio InnerContentAudioStruct `json:"input_audio"`
}

type InnerContentImageStruct struct {
	Detail string `json:"detail"`
	Url    string `json:"url"`
}

type InnerContentAudioStruct struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

func ToAIContent(chunk rawopenai.ChatCompletionChunkChoice) []contents.IAIContent {
	var result []contents.IAIContent
	var strSlice []string

	if err := json.Unmarshal([]byte(chunk.Delta.Content), &strSlice); err != nil {
		return []contents.IAIContent{contents.NewTextContentWithRefusal(chunk.Delta.Content, chunk.Delta.Refusal)}
	}

	for _, input := range strSlice {
		result = append(result, parseContent(input, chunk.Delta.Refusal))
	}

	return result
}

func parseContent(input string, refusal string) contents.IAIContent {
	var jsonObj InnerContentStruct
	if err := json.Unmarshal([]byte(input), &jsonObj); err != nil {
		return contents.NewTextContentWithRefusal(input, refusal)
	}

	switch jsonObj.Type {
	case "image":
		return contents.NewDataContentWithRefusal(jsonObj.Image.Url, "image", refusal)
	default:
		return contents.NewTextContentWithRefusal(input, refusal)
	}
}

func ToOpenAIEmbeddingParams[TInput any](values []TInput, options *embeddings.EmbeddingGenerationOptions) *rawopenai.EmbeddingNewParams {
	return nil
}

func ToGeneratedEmbeddings(res *rawopenai.CreateEmbeddingResponse) *embeddings.GeneratedEmbeddings[embeddings.EmbeddingT[float64]] {
	return nil
}

func ToChatResponseUpdateFromAssistantStreamEvent(evt rawopenai.AssistantStreamEvent) *chatcompletion.ChatResponseUpdate {
	panic("unimplemented")
}
