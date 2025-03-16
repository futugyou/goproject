package openai

import (
	"encoding/json"

	"github.com/futugyou/yomawari/generative-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/generative-ai/abstractions/contents"
	"github.com/futugyou/yomawari/generative-ai/abstractions/embeddings"
	rawopenai "github.com/openai/openai-go"
)

// Conversion functions
func ToOpenAIChatCompletion(response chatcompletion.ChatResponse, options json.RawMessage) rawopenai.ChatCompletion {
	result := rawopenai.ChatCompletion{}
	return result
}

func ToOpenAIMessages(chatMessages []chatcompletion.ChatMessage) []rawopenai.ChatCompletionMessageParamUnion {
	return nil
}

func ToOpenAIChatRequest(options *chatcompletion.ChatOptions) *rawopenai.ChatCompletionNewParams {
	return nil
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
	return nil
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
