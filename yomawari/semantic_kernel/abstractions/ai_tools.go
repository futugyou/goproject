package abstractions

import (
	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	aicontents "github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	aifunctions "github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
)

func ToChatMessageContent(message chatcompletion.ChatMessage, response *chatcompletion.ChatResponse) ChatMessageContent {
	result := &ChatMessageContent{
		Metadata: message.AdditionalProperties,
		Role:     CreateAuthorRole(string(message.Role)),
	}
	if message.AuthorName != nil {
		result.AuthorName = *message.AuthorName
	}
	if response != nil {
		if response.ModelId != nil {
			result.ModelId = *response.ModelId
		}
		if response.RawRepresentation != nil {
			result.InnerContent = response.RawRepresentation
		} else {
			result.InnerContent = message.RawRepresentation
		}
	}

	for _, content := range message.Contents {
		var resultContent KernelContent
		switch tc := content.(type) {
		case aicontents.TextContent:
			resultContent = &TextContent{
				ModelId:      result.ModelId,
				Metadata:     content.GetAdditionalProperties(),
				InnerContent: content.GetRawRepresentation(),
				Text:         tc.Text,
			}
		case aicontents.DataContent:
			if tc.MediaTypeStartsWith("image") {
				resultContent = &ImageContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			} else if tc.MediaTypeStartsWith("audio") {
				resultContent = &AudioContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			} else {
				resultContent = &BinaryContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			}
		case aicontents.UriContent:
			if tc.MediaTypeStartsWith("image") {
				resultContent = &ImageContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			} else if tc.MediaTypeStartsWith("audio") {
				resultContent = &AudioContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			} else {
				resultContent = &BinaryContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			}
		case aicontents.FunctionCallContent:
			resultContent = &FunctionCallContent{
				ModelId:      result.ModelId,
				Metadata:     content.GetAdditionalProperties(),
				Id:           tc.CallId,
				FunctionName: tc.Name,
				Arguments:    tc.Arguments,
				InnerContent: content.GetRawRepresentation(),
			}
		case aicontents.FunctionResultContent:
			resultContent = &FunctionResultContent{
				ModelId:      result.ModelId,
				Metadata:     content.GetAdditionalProperties(),
				CallId:       tc.CallId,
				Result:       tc.Result,
				InnerContent: content.GetRawRepresentation(),
			}
		}

		if resultContent != nil {
			result.Items.Add(resultContent)
		}
	}

	return *result
}

func ToPromptExecutionSettings(options *chatcompletion.ChatOptions) *PromptExecutionSettings {
	if options == nil {
		return nil
	}

	settings := &PromptExecutionSettings{
		ExtensionData: make(map[string]interface{}),
	}

	if options.ModelId != nil {
		settings.ModelId = *options.ModelId
	}
	if options.Temperature != nil {
		settings.ExtensionData["temperature"] = *options.Temperature
	}
	if options.MaxOutputTokens != nil {
		settings.ExtensionData["max_tokens"] = *options.MaxOutputTokens
	}
	if options.FrequencyPenalty != nil {
		settings.ExtensionData["frequency_penalty"] = *options.FrequencyPenalty
	}
	if options.PresencePenalty != nil {
		settings.ExtensionData["presence_penalty"] = *options.PresencePenalty
	}
	if options.StopSequences != nil {
		settings.ExtensionData["stop_sequences"] = options.StopSequences
	}
	if options.TopP != nil {
		settings.ExtensionData["top_p"] = *options.TopP
	}
	if options.TopK != nil {
		settings.ExtensionData["top_k"] = *options.TopK
	}
	if options.Seed != nil {
		settings.ExtensionData["seed"] = *options.Seed
	}

	if options.ResponseFormat != nil {
		settings.ExtensionData["response_format"] = *options.Seed
	}

	if options.AdditionalProperties != nil {
		for k, v := range options.AdditionalProperties {
			if v != nil {
				settings.ExtensionData[k] = v
			}
		}
	}

	if len(options.Tools) > 0 {
		var fs []KernelFunction
		for _, tool := range options.Tools {
			if fn, ok := tool.(aifunctions.AIFunction); ok {
				fs = append(fs, AIFunctionTosKernelFunction(fn))
			}
		}

		tm := *options.ToolMode
		if tm == chatcompletion.AutoMode {
			settings.FunctionChoiceBehavior = *NewFunctionChoiceBehavior(fs, nil, "auto", false)
		}
		if tm == chatcompletion.RequireAnyMode {
			// TODO use struct to reflect chatcompletion.ChatToolMode, need RequiredFunctionName field
			var requiredFunctionName *string = getRequiredFunctionName() //	tm.RequiredFunctionName
			if requiredFunctionName == nil {
				settings.FunctionChoiceBehavior = *NewFunctionChoiceBehavior(fs, nil, "required", false)
			} else {
				var matched []KernelFunction
				for _, f := range fs {
					if f.GetName() == *requiredFunctionName {
						matched = append(matched, f)
					}
				}
				settings.FunctionChoiceBehavior = *NewFunctionChoiceBehavior(matched, nil, "required", false)
			}
		}
	}

	return settings
}

func getRequiredFunctionName() *string {
	return nil
}

func ToChatHistory(chatMessages []chatcompletion.ChatMessage) ChatHistory {
	c := &ChatHistory{List: *core.NewList[ChatMessageContent]()}
	for _, v := range chatMessages {
		c.Add(ToChatMessageContent(v, nil))
	}
	return *c
}
