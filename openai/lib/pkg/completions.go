package pkg

const completionsPath string = "completions"

type CreateCompletionRequest struct {
	Model            string          `json:"model,omitempty"`
	Prompt           interface{}     `json:"prompt,omitempty"`
	Suffix           string          `json:"suffix,omitempty"`
	MaxTokens        int32           `json:"max_tokens,omitempty"`
	Temperature      float32         `json:"temperature,omitempty"`
	Top_p            float32         `json:"top_p,omitempty"`
	N                int32           `json:"n,omitempty"`
	Stream           bool            `json:"stream,omitempty"`
	Logprobs         int32           `json:"logprobs,omitempty"`
	Echo             bool            `json:"echo,omitempty"`
	Stop             interface{}     `json:"stop,omitempty"`
	PresencePenalty  float32         `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32         `json:"frequency_penalty,omitempty"`
	BestOf           int32           `json:"best_of,omitempty"`
	LogitBias        map[string]int8 `json:"logit_bias,omitempty"`
	User             string          `json:"user,omitempty"`
}

type CreateCompletionResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	ID      string       `json:"id,omitempty"`
	Object  string       `json:"object,omitempty"`
	Created int32        `json:"created,omitempty"`
	Model   string       `json:"model,omitempty"`
	Choices []choices    `json:"choices,omitempty"`
	Usage   *usage       `json:"usage,omitempty"`
}

type choices struct {
	Text         string `json:"text"`
	Index        int32  `json:"index"`
	Logprobs     int32  `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type usage struct {
	PromptTokens     int32 `json:"prompt_tokens"`
	CompletionTokens int32 `json:"completion_tokens"`
	TotalTokens      int32 `json:"total_tokens"`
}

func (client *openaiClient) CreateCompletion(request CreateCompletionRequest) *CreateCompletionResponse {
	result := &CreateCompletionResponse{}
	client.Post(completionsPath, request, result)
	return result
}
