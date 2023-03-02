package lib

const chatCompletionPath string = "chat/completions"

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
	client.Post(chatCompletionPath, request, result)
	return result
}
