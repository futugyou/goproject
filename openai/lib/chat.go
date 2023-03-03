package lib

import (
	"fmt"

	"golang.org/x/exp/slices"
)

const chatCompletionPath string = "chat/completions"

var supportedChatModel = []string{GPT35_turbo, GPT35_turbo_0301}

var chatModelError = func(model string) *OpenaiError {
	return &OpenaiError{
		ErrorMessage: "Currently, only gpt-3.5-turbo and gpt-3.5-turbo-0301 are supported.",
		ErrorType:    "invalid parameters",
		Param:        fmt.Sprintf("current chat model is: %s", model),
	}
}

type CreateChatCompletionRequest struct {
	Model            string                `json:"model"`
	Messages         ChatCompletionMessage `json:"messages"`
	Temperature      float32               `json:"temperature,omitempty"`
	Top_p            float32               `json:"top_p,omitempty"`
	N                int32                 `json:"n,omitempty"`
	Stream           bool                  `json:"stream,omitempty"`
	Stop             []string              `json:"stop,omitempty"`
	MaxTokens        int32                 `json:"max_tokens,omitempty"`
	PresencePenalty  float32               `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32               `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int32      `json:"logit_bias,omitempty"`
	User             string                `json:"user,omitempty"`
}

type CreateChatCompletionResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	ID      string       `json:"id,omitempty"`
	Object  string       `json:"object,omitempty"`
	Created int32        `json:"created,omitempty"`
	Model   string       `json:"model,omitempty"`
	Choices []Choices    `json:"choices,omitempty"`
	Usage   *Usage       `json:"usage,omitempty"`
}

func (client *openaiClient) CreateChatCompletion(request CreateChatCompletionRequest) *CreateChatCompletionResponse {
	result := &CreateChatCompletionResponse{}

	err := validateChatModel(request.Model)
	if err != nil {
		result.Error = err
		return result
	}

	client.Post(chatCompletionPath, request, result)
	return result
}

func validateChatModel(model string) *OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedChatModel, model) {
		return chatModelError(model)
	}

	return nil
}
