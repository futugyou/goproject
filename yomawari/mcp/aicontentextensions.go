package mcp

import (
	"encoding/base64"
	"encoding/json"

	"github.com/futugyou/yomawari/extensions-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions-ai/abstractions/contents"
	"github.com/futugyou/yomawari/mcp/protocol/types"
)

func ResourceContentsToAIContent(content types.IResourceContents) contents.IAIContent {
	var c contents.IAIContent
	switch content := content.(type) {
	case *types.BlobResourceContents:
		decoded, err := base64.URLEncoding.DecodeString(content.Blob)
		if err != nil {
			mimeType := "application/octet-stream"
			if content.MimeType != nil && len(*content.MimeType) > 0 {
				mimeType = *content.MimeType
			}
			d := contents.NewDataContent(string(decoded), mimeType)
			d.AddAdditionalProperty("uri", content.Uri)
			c = d
		}

	case *types.TextResourceContents:
		d := contents.NewTextContent(content.Text)
		d.AddAdditionalProperty("uri", content.Uri)
		c = d
	}
	return c
}

func ContentToAIContent(content types.Content) contents.IAIContent {
	var c contents.IAIContent

	if (content.Type == "image" || content.Type == "audio") && content.MimeType != nil && content.Data != nil {
		decoded, err := base64.URLEncoding.DecodeString(*content.Data)
		if err != nil {
			d := contents.NewDataContent(string(decoded), *content.MimeType)
			d.RawRepresentation = content
			c = d
		}
	} else if content.Type == "resource" || content.Resource != nil {
		c = ResourceContentsToAIContent(content.Resource)
	} else {
		d := contents.NewTextContent(*content.Text)
		d.RawRepresentation = content
		c = d
	}

	return c
}

func ChatMessageToPromptMessages(chatMessage chatcompletion.ChatMessage) []types.PromptMessage {
	r := types.RoleAssistant
	if chatMessage.Role == chatcompletion.RoleUser {
		r = types.RoleUser
	}
	messages := []types.PromptMessage{}

	for _, content := range chatMessage.Contents {
		if c, ok := content.(contents.TextContent); ok {
			messages = append(messages, types.PromptMessage{Role: r, Content: AIContentToContent(c)})
		}
		if c, ok := content.(contents.DataContent); ok {
			messages = append(messages, types.PromptMessage{Role: r, Content: AIContentToContent(c)})
		}
	}
	return messages
}

func AIContentToContent(content contents.IAIContent) types.Content {
	switch content := content.(type) {
	case *contents.TextContent:
		return types.Content{
			Type: "text",
			Text: &content.Text,
		}
	case *contents.DataContent:
		c := types.Content{
			Type:     "resource",
			MimeType: &content.MediaType,
		}
		decoded := base64.URLEncoding.EncodeToString(content.Data)
		c.Data = &decoded
		if content.MediaTypeStartsWith("image") {
			c.Type = "image"
		} else if content.MediaTypeStartsWith("audio") {
			c.Type = "audio"
		}
		return c
	default:
		data, err := json.Marshal(content.(contents.AIContent))
		if err != nil {
			data = []byte{}
		}
		d := string(data)
		return types.Content{
			Type: "text",
			Text: &d,
		}
	}
}

func ResourceContentsListToAIContents(cont []types.IResourceContents) []contents.IAIContent {
	list := []contents.IAIContent{}
	for _, content := range cont {
		list = append(list, ResourceContentsToAIContent(content))
	}
	return list
}

func ListContentToAIContents(cont []types.Content) []contents.IAIContent {
	list := []contents.IAIContent{}
	for _, content := range cont {
		list = append(list, ContentToAIContent(content))
	}
	return list
}

func ToChatMessages(promptResult types.GetPromptResult) []chatcompletion.ChatMessage {
	list := []chatcompletion.ChatMessage{}
	for _, v := range promptResult.Messages {
		list = append(list, PromptMessageToChatMessage(v))
	}
	return list
}

func PromptMessageToChatMessage(promptMessage types.PromptMessage) chatcompletion.ChatMessage {
	role := chatcompletion.RoleAssistant
	if promptMessage.Role == types.RoleUser {
		role = chatcompletion.RoleUser
	}
	return chatcompletion.ChatMessage{
		RawRepresentation: promptMessage,
		Role:              role,
		Contents:          []contents.IAIContent{ContentToAIContent(promptMessage.Content)},
	}
}
