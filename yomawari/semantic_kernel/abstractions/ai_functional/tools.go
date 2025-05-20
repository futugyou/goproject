package ai_functional

import (
	"encoding/json"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	aicontents "github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	aifunctions "github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/contents"
)

func ToChatMessage(content contents.ChatMessageContent) chatcompletion.ChatMessage {
	message := &chatcompletion.ChatMessage{
		AdditionalProperties: content.Metadata,
		AuthorName:           &content.AuthorName,
		RawRepresentation:    content.InnerContent,
	}

	message.Role = chatcompletion.StringToChatRole(string(content.Role))

	for _, item := range content.Items.Items() {
		var aiContent aicontents.IAIContent

		switch tc := item.(type) {
		case contents.TextContent:
			aiContent = aicontents.NewTextContent(tc.Text)
		case contents.ImageContent:
			if len(tc.DataUri) > 0 {
				aiContent = aicontents.NewDataContent(tc.DataUri, tc.MimeType)
			}

			if len(tc.Uri.String()) > 0 {
				mimeType := "image/*"
				if len(tc.MimeType) > 0 {
					mimeType = tc.MimeType
				}
				aiContent = aicontents.UriContent{URI: tc.Uri.String(), MediaType: mimeType}
			}

		case contents.AudioContent:
			if len(tc.DataUri) > 0 {
				aiContent = aicontents.NewDataContent(tc.DataUri, tc.MimeType)
			}

			if len(tc.Uri.String()) > 0 {
				mimeType := "audio/*"
				if len(tc.MimeType) > 0 {
					mimeType = tc.MimeType
				}
				aiContent = aicontents.UriContent{URI: tc.Uri.String(), MediaType: mimeType}
			}
		case contents.BinaryContent:
			if len(tc.DataUri) > 0 {
				aiContent = aicontents.NewDataContent(tc.DataUri, tc.MimeType)
			}

			if len(tc.Uri.String()) > 0 {
				mimeType := "application/octet-stream"
				if len(tc.MimeType) > 0 {
					mimeType = tc.MimeType
				}
				aiContent = aicontents.UriContent{URI: tc.Uri.String(), MediaType: mimeType}
			}
		case contents.FunctionCallContent:
			aiContent = aicontents.FunctionCallContent{CallId: tc.Id, Name: tc.FunctionName, Arguments: tc.Arguments}
		case contents.FunctionResultContent:
			aiContent = aicontents.FunctionResultContent{CallId: tc.CallId, Result: tc.Result}
		}

		if aiContent != nil {
			message.Contents = append(message.Contents, aiContent)
		}
	}

	return *message
}

func ToChatMessageList(chatHistory ChatHistory) []chatcompletion.ChatMessage {
	result := []chatcompletion.ChatMessage{}
	for _, v := range chatHistory.Items() {
		result = append(result, ToChatMessage(v))
	}
	return result
}

func ToChatMessageContent(message chatcompletion.ChatMessage, response *chatcompletion.ChatResponse) contents.ChatMessageContent {
	result := &contents.ChatMessageContent{
		Metadata: message.AdditionalProperties,
		Role:     contents.CreateAuthorRole(string(message.Role)),
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
		var resultContent contents.KernelContent
		switch tc := content.(type) {
		case aicontents.TextContent:
			resultContent = &contents.TextContent{
				ModelId:      result.ModelId,
				Metadata:     content.GetAdditionalProperties(),
				InnerContent: content.GetRawRepresentation(),
				Text:         tc.Text,
			}
		case aicontents.DataContent:
			if tc.MediaTypeStartsWith("image") {
				resultContent = &contents.ImageContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			} else if tc.MediaTypeStartsWith("audio") {
				resultContent = &contents.AudioContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			} else {
				resultContent = &contents.BinaryContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			}
		case aicontents.UriContent:
			if tc.MediaTypeStartsWith("image") {
				resultContent = &contents.ImageContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			} else if tc.MediaTypeStartsWith("audio") {
				resultContent = &contents.AudioContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			} else {
				resultContent = &contents.BinaryContent{
					ModelId:      result.ModelId,
					Metadata:     content.GetAdditionalProperties(),
					InnerContent: content.GetRawRepresentation(),
					DataUri:      tc.URI,
					MimeType:     tc.MediaType,
				}
			}
		case aicontents.FunctionCallContent:
			resultContent = &contents.FunctionCallContent{
				ModelId:      result.ModelId,
				Metadata:     content.GetAdditionalProperties(),
				Id:           tc.CallId,
				FunctionName: tc.Name,
				Arguments:    tc.Arguments,
				InnerContent: content.GetRawRepresentation(),
			}
		case aicontents.FunctionResultContent:
			resultContent = &contents.FunctionResultContent{
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

func ToStreamingChatCompletionUpdate(content contents.StreamingChatMessageContent) *chatcompletion.ChatResponseUpdate {
	r := chatcompletion.StringToChatRole(string(content.Role))
	update := &chatcompletion.ChatResponseUpdate{
		AdditionalProperties: content.Metadata,
		AuthorName:           &content.AuthorName,
		ModelId:              &content.ModelId,
		RawRepresentation:    content,
		Role:                 &r,
	}

	for _, item := range content.Items.Items() {
		var aiContent aicontents.IAIContent
		switch tc := item.(type) {
		case contents.StreamingTextContent:
			aiContent = aicontents.NewTextContent(tc.Text)
		case contents.StreamingFunctionCallUpdateContent:
			var a map[string]interface{}
			json.Unmarshal([]byte(tc.Arguments), &a)
			aiContent = &aicontents.FunctionCallContent{
				AIContent: &aicontents.AIContent{},
				CallId:    tc.CallId,
				Name:      tc.Name,
				Arguments: a,
			}
		}

		if aiContent != nil {
			update.Contents = append(update.Contents, aiContent)
		}
	}

	return update
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
		// var fs []functions.KernelFunction
		for _, tool := range options.Tools {
			if _, ok := tool.(aifunctions.AIFunction); ok {
				// fs = append(fs, fn.AsKernelFunction())
			}
		}

		tm := *options.ToolMode
		if tm == chatcompletion.AutoMode {
			// settings.FunctionChoiceBehavior = FunctionChoiceBehaviorAuto(fs, false)
		}
		if tm == chatcompletion.RequireAnyMode {
			// if tm.RequiredFunctionName == nil {
			// 	settings.FunctionChoiceBehavior = FunctionChoiceBehaviorRequired(fs, false)
			// } else {
			// var matched []functions.KernelFunction
			// for _, f := range fs {
			// 	if f.Name == *tm.RequiredFunctionName {
			// 		matched = append(matched, f)
			// 	}
			// }
			// settings.FunctionChoiceBehavior = FunctionChoiceBehaviorRequired(matched, false)
			// }
		}
	}

	return settings
}
