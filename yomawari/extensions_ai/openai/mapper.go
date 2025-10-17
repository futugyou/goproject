package openai

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/futugyou/yomawari/extensions_ai/abstractions"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/embeddings"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
	rawopenai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
	"github.com/openai/openai-go/v3/shared/constant"
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
	result := rawopenai.ChatCompletionMessageParamUnion{}

	name := ""
	if chatMessage.AuthorName != nil {
		name = *chatMessage.AuthorName
	}

	var role string = (string)(chatMessage.Role)
	conts := ToOpenAIContents(chatMessage.Contents, role)
	switch role {
	case "developer":
		content := rawopenai.ChatCompletionDeveloperMessageParamContentUnion{}
		for i := 0; i < len(conts); i++ {
			content.OfArrayOfContentParts = append(content.OfArrayOfContentParts, rawopenai.ChatCompletionContentPartTextParam{Text: conts[i].Content})
		}
		result.OfDeveloper = &rawopenai.ChatCompletionDeveloperMessageParam{
			Name:    param.NewOpt(name),
			Content: content,
		}
	case "system":
		content := rawopenai.ChatCompletionSystemMessageParamContentUnion{}
		for i := 0; i < len(conts); i++ {
			content.OfArrayOfContentParts = append(content.OfArrayOfContentParts, rawopenai.ChatCompletionContentPartTextParam{Text: conts[i].Content})
		}
		result.OfSystem = &rawopenai.ChatCompletionSystemMessageParam{
			Name:    param.NewOpt(name),
			Content: content,
		}
	case "assistant":
		content := rawopenai.ChatCompletionAssistantMessageParamContentUnion{}
		for i := 0; i < len(conts); i++ {
			content.OfArrayOfContentParts = append(content.OfArrayOfContentParts,
				rawopenai.ChatCompletionAssistantMessageParamContentArrayOfContentPartUnion{
					OfText: &rawopenai.ChatCompletionContentPartTextParam{
						Text: conts[i].Content,
					},
				},
			)
		}
		result.OfAssistant = &rawopenai.ChatCompletionAssistantMessageParam{
			Name:    param.NewOpt(name),
			Content: content,
		}
	case "tool":
		content := rawopenai.ChatCompletionToolMessageParamContentUnion{}
		for i := 0; i < len(conts); i++ {
			content.OfArrayOfContentParts = append(content.OfArrayOfContentParts, rawopenai.ChatCompletionContentPartTextParam{Text: conts[i].Content})
		}
		toolid := ""
		for _, con := range chatMessage.Contents {
			if c, ok := con.(*contents.FunctionCallContent); ok {
				if len(c.CallId) > 0 {
					toolid = c.CallId
					break
				}
			}
		}
		result.OfTool = &rawopenai.ChatCompletionToolMessageParam{
			ToolCallID: toolid,
			Content:    content,
		}
	default:
		content := rawopenai.ChatCompletionUserMessageParamContentUnion{}
		for i := 0; i < len(conts); i++ {
			switch conts[i].Type {
			case "image":
				content.OfArrayOfContentParts = append(content.OfArrayOfContentParts,
					rawopenai.ChatCompletionContentPartUnionParam{
						OfImageURL: &rawopenai.ChatCompletionContentPartImageParam{
							ImageURL: rawopenai.ChatCompletionContentPartImageImageURLParam{
								URL:    conts[i].Content,
								Detail: "auto",
							},
						},
					},
				)
			case "audio":
				content.OfArrayOfContentParts = append(content.OfArrayOfContentParts,
					rawopenai.ChatCompletionContentPartUnionParam{
						OfInputAudio: &rawopenai.ChatCompletionContentPartInputAudioParam{
							InputAudio: rawopenai.ChatCompletionContentPartInputAudioInputAudioParam{
								Data:   conts[i].Content,
								Format: conts[i].Format,
							},
						},
					},
				)
			default:
				content.OfArrayOfContentParts = append(content.OfArrayOfContentParts,
					rawopenai.ChatCompletionContentPartUnionParam{
						OfText: &rawopenai.ChatCompletionContentPartTextParam{
							Text: conts[i].Content},
					},
				)
			}

		}
		result.OfUser = &rawopenai.ChatCompletionUserMessageParam{
			Name:    param.NewOpt(name),
			Content: content,
		}

	}

	return result
}

type ContentInfo struct {
	Content string
	Type    string
	Format  string
}

func ToOpenAIContents(cons []contents.IAIContent, role string) []ContentInfo {
	result := []ContentInfo{}
	for _, v := range cons {
		switch con := v.(type) {
		case *contents.TextContent:
			result = append(result, ContentInfo{Content: con.Text, Type: "text"})
		case *contents.UriContent:
			if con.MediaTypeStartsWith("image") {
				result = append(result, ContentInfo{Content: con.URI, Type: "image"})
			}
		case *contents.DataContent:
			if con.MediaTypeStartsWith("image") {
				result = append(result, ContentInfo{Content: con.GetURI(), Type: "image"})
			} else if con.MediaTypeStartsWith("audio") {
				format := "mp3"
				if con.MediaTypeStartsWith("audio/wav") {
					format = "wav"
				}
				result = append(result, ContentInfo{Content: con.GetURI(), Type: "audio", Format: format})
			} else {
				result = append(result, ContentInfo{Content: string(con.Data), Type: "text"})
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
		result.Model = shared.ChatModel(*options.ModelId)
	}
	if options.FrequencyPenalty != nil {
		result.FrequencyPenalty = param.NewOpt(*options.FrequencyPenalty)
	}
	if options.MaxOutputTokens != nil {
		result.MaxCompletionTokens = param.NewOpt(*options.MaxOutputTokens)
	}
	if options.PresencePenalty != nil {
		result.PresencePenalty = param.NewOpt(*options.PresencePenalty)
	}
	if options.Seed != nil {
		result.Seed = param.NewOpt(*options.Seed)
	}
	if options.Temperature != nil {
		result.Temperature = param.NewOpt(*options.Temperature)
	}
	if options.TopP != nil {
		result.TopP = param.NewOpt(*options.TopP)
	}
	if options.AllowMultipleToolCalls != nil {
		result.ParallelToolCalls = param.NewOpt(*options.AllowMultipleToolCalls)
	}

	if v, ok := options.AdditionalProperties["User"].(string); ok {
		result.User = param.NewOpt(v)
	}

	if v, ok := options.AdditionalProperties["ReasoningEffort"].(string); ok {
		result.ReasoningEffort = shared.ReasoningEffort(v)
	}
	if v, ok := options.AdditionalProperties["Modalities"].([]string); ok {
		result.Modalities = v
	}
	if v, ok := options.AdditionalProperties["Store"].(bool); ok {
		result.Store = param.NewOpt(v)
	}
	if v, ok := options.AdditionalProperties["TopLogprobs"].(int64); ok {
		result.TopLogprobs = param.NewOpt(v)
	}
	return result
}

func ToOpenAIChatCompletionToolUnionParam(v functions.AIFunction) rawopenai.ChatCompletionToolUnionParam {
	var m shared.FunctionParameters = v.GetJsonSchema()
	strict := false
	if v, ok := v.GetAdditionalProperties()["strictJsonSchema"].(bool); ok {
		strict = v
	}
	pa := &rawopenai.ChatCompletionFunctionToolParam{
		Function: shared.FunctionDefinitionParam{
			Name:        v.GetName(),
			Description: param.NewOpt(v.GetDescription()),
			Strict:      param.NewOpt(strict),
			Parameters:  m,
		}}

	p := rawopenai.ChatCompletionToolUnionParam{
		OfFunction: pa,
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
		con := contents.NewTextContent(v.Content)
		message.Contents = append(message.Contents, con)
	}

	if len(v.Audio.Data) > 0 {
		con := contents.DataContent{
			AIContent: contents.NewAIContent(nil, nil),
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

func ToChatRole(v constant.Assistant) chatcompletion.ChatRole {
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
			AIContent: contents.NewAIContent(nil, nil),
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

type ToolCallsCache struct {
	sync.Mutex
	data map[string]rawopenai.ChatCompletionChunkChoiceDeltaToolCall
}

func ToChatResponseUpdateWithFunctions(response *rawopenai.ChatCompletionChunk, toolCallsCache *ToolCallsCache) *chatcompletion.ChatResponseUpdate {
	if response == nil || len(response.Choices) == 0 {
		return nil
	}

	toolCallsCache.Mutex.Lock()
	defer toolCallsCache.Mutex.Unlock()

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

	finishReason := chatcompletion.ChatFinishReason((string)(response.Choices[len(response.Choices)-1].FinishReason))
	role := chatcompletion.StringToChatRole((string)(response.Choices[len(response.Choices)-1].Delta.Role))

	result.Role = &role
	result.FinishReason = &finishReason
	refusal := ""
	for _, choice := range response.Choices {
		if len(choice.Delta.Refusal) > 0 {
			refusal = choice.Delta.Refusal
		}
		for _, toolCall := range choice.Delta.ToolCalls {
			existing, found := toolCallsCache.data[toolCall.ID]
			if found {
				if toolCall.Function.Name != "" {
					existing.Function.Name = toolCall.Function.Name
				}
				if toolCall.Function.Arguments != "" {
					existing.Function.Arguments += toolCall.Function.Arguments
				}
				toolCallsCache.data[toolCall.ID] = existing
			} else {
				toolCallsCache.data[toolCall.ID] = toolCall
			}
		}
	}

	for _, v := range toolCallsCache.data {
		con := contents.CreateFromParsedArguments(v.Function.Arguments, v.ID, v.Function.Name, func(args string) (map[string]interface{}, error) {
			var result map[string]interface{}
			err := json.Unmarshal([]byte(args), &result)
			return result, err
		})
		result.Contents = append(result.Contents, con)
	}

	if len(refusal) > 0 {
		result.AdditionalProperties["refusal"] = refusal
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
		return []contents.IAIContent{contents.NewTextContent(chunk.Delta.Content)}
	}

	for _, input := range strSlice {
		result = append(result, parseContent(input))
	}

	return result
}

func parseContent(input string) contents.IAIContent {
	var jsonObj InnerContentStruct
	if err := json.Unmarshal([]byte(input), &jsonObj); err != nil {
		return contents.NewTextContent(input)
	}

	switch jsonObj.Type {
	case "image":
		return contents.NewDataContent(jsonObj.Image.Url, "image")
	default:
		return contents.NewTextContent(input)
	}
}

func ToOpenAIEmbeddingParams(values []string, options *embeddings.EmbeddingGenerationOptions) *rawopenai.EmbeddingNewParams {
	if options == nil || len(values) == 0 {
		return nil
	}
	var i rawopenai.EmbeddingNewParamsInputUnion = rawopenai.EmbeddingNewParamsInputUnion{OfArrayOfStrings: values}
	result := &rawopenai.EmbeddingNewParams{
		Input:      i,
		Model:      rawopenai.EmbeddingModel(*options.ModelId),
		Dimensions: param.NewOpt(*options.Dimensions),
	}

	if v, ok := options.AdditionalProperties["encoding_format"].(string); ok {
		if v == "base64" {
			result.EncodingFormat = "base64"
		} else {
			result.EncodingFormat = "float"
		}

	}
	if v, ok := options.AdditionalProperties["user"].(string); ok {
		result.User = param.NewOpt(v)
	}
	return result
}

func ToGeneratedEmbeddings(res *rawopenai.CreateEmbeddingResponse) *embeddings.GeneratedEmbeddings[embeddings.EmbeddingT[float64]] {
	if res == nil {
		return nil
	}
	emb := []embeddings.EmbeddingT[float64]{}
	t := time.Now().UTC()
	for _, v := range res.Data {
		emb = append(emb, embeddings.EmbeddingT[float64]{
			Embedding: embeddings.Embedding{
				CreatedAt:            &t,
				ModelId:              &res.Model,
				AdditionalProperties: map[string]interface{}{},
			},
			Vector: v.Embedding,
		})
	}
	result := embeddings.NewGeneratedEmbeddingsFromCollection(emb)
	result.Usage = &abstractions.UsageDetails{
		InputTokenCount: &res.Usage.PromptTokens,
		TotalTokenCount: &res.Usage.TotalTokens,
	}
	return result
}

func ConvertContent(con rawopenai.MessageContentUnion) contents.IAIContent {
	if len(con.ImageURL.URL) > 0 {
		return contents.NewDataContentWithRefusal(con.ImageURL.URL, "image", con.Refusal)
	}

	if len(con.ImageFile.FileID) > 0 {
		return contents.NewTextContentWithRefusal(con.ImageFile.FileID, con.Refusal)
	}

	if len(con.Text.Value) > 0 {
		return contents.NewTextContentWithRefusal(con.Text.Value, con.Refusal)
	}

	return nil
}

func GetUsageContentStep(runUsage rawopenai.RunStepUsage) contents.UsageContent {
	return contents.UsageContent{
		AIContent: contents.NewAIContent(nil, nil),
		Details: abstractions.UsageDetails{
			InputTokenCount:      &runUsage.PromptTokens,
			OutputTokenCount:     &runUsage.CompletionTokens,
			TotalTokenCount:      &runUsage.TotalTokens,
			AdditionalProperties: map[string]int64{},
		},
	}
}

func GetUsageContent(runUsage rawopenai.RunUsage) contents.UsageContent {
	return contents.UsageContent{
		AIContent: contents.NewAIContent(nil, nil),
		Details: abstractions.UsageDetails{
			InputTokenCount:      &runUsage.PromptTokens,
			OutputTokenCount:     &runUsage.CompletionTokens,
			TotalTokenCount:      &runUsage.TotalTokens,
			AdditionalProperties: map[string]int64{},
		},
	}
}
