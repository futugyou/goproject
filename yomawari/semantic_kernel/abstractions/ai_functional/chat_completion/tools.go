package chat_completion

import (
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	aicontents "github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
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
