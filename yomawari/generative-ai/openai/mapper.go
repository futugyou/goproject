package openai

import (
	"encoding/json"

	"github.com/futugyou/yomawari/generative-ai/abstractions/chatcompletion"
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

func ToChatResponse(response *rawopenai.ChatCompletion) *chatcompletion.ChatResponse {
	return nil
}

func ToChatResponseUpdate(response *rawopenai.ChatCompletionChunk) *chatcompletion.ChatResponseUpdate {
	return nil
}
