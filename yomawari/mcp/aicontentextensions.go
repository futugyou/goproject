package mcp

import (
	"encoding/base64"
	"encoding/json"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/futugyou/yomawari/mcp/protocol"
)

func ResourceContentsToAIContent(content protocol.IResourceContents) contents.IAIContent {
	var c contents.IAIContent
	switch content := content.(type) {
	case *protocol.BlobResourceContents:
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

	case *protocol.TextResourceContents:
		d := contents.NewTextContent(content.Text)
		d.AddAdditionalProperty("uri", content.Uri)
		c = d
	}
	return c
}

func ContentToAIContent(content protocol.Content) contents.IAIContent {
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

func ChatMessageToPromptMessages(chatMessage chatcompletion.ChatMessage) []protocol.PromptMessage {
	r := protocol.RoleAssistant
	if chatMessage.Role == chatcompletion.RoleUser {
		r = protocol.RoleUser
	}
	messages := []protocol.PromptMessage{}

	for _, content := range chatMessage.Contents {
		if c, ok := content.(contents.TextContent); ok {
			messages = append(messages, protocol.PromptMessage{Role: r, Content: AIContentToContent(c)})
		}
		if c, ok := content.(contents.DataContent); ok {
			messages = append(messages, protocol.PromptMessage{Role: r, Content: AIContentToContent(c)})
		}
	}
	return messages
}

func AIContentToContent(content contents.IAIContent) protocol.Content {
	switch content := content.(type) {
	case *contents.TextContent:
		return protocol.Content{
			Type: "text",
			Text: &content.Text,
		}
	case *contents.DataContent:
		c := protocol.Content{
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
		data, err := json.Marshal(content.(*contents.AIContent))
		if err != nil {
			data = []byte{}
		}
		d := string(data)
		return protocol.Content{
			Type: "text",
			Text: &d,
		}
	}
}

func ResourceContentsListToAIContents(cont []protocol.IResourceContents) []contents.IAIContent {
	list := []contents.IAIContent{}
	for _, content := range cont {
		list = append(list, ResourceContentsToAIContent(content))
	}
	return list
}

func ListContentToAIContents(cont []protocol.Content) []contents.IAIContent {
	list := []contents.IAIContent{}
	for _, content := range cont {
		list = append(list, ContentToAIContent(content))
	}
	return list
}

func ToChatMessages(promptResult protocol.GetPromptResult) []chatcompletion.ChatMessage {
	list := []chatcompletion.ChatMessage{}
	for _, v := range promptResult.Messages {
		list = append(list, PromptMessageToChatMessage(v))
	}
	return list
}

func PromptMessageToChatMessage(promptMessage protocol.PromptMessage) chatcompletion.ChatMessage {
	role := chatcompletion.RoleAssistant
	if promptMessage.Role == protocol.RoleUser {
		role = chatcompletion.RoleUser
	}
	return chatcompletion.ChatMessage{
		RawRepresentation: promptMessage,
		Role:              role,
		Contents:          []contents.IAIContent{ContentToAIContent(promptMessage.Content)},
	}
}
