package lib

type Choices struct {
	Text         string                  `json:"text"`
	Index        int32                   `json:"index"`
	Logprobs     *Logprobs               `json:"logprobs,omitempty"`
	FinishReason string                  `json:"finish_reason,omitempty"`
	Message      []ChatCompletionMessage `json:"message,omitempty"`
}

type Usage struct {
	PromptTokens     int32 `json:"prompt_tokens"`
	CompletionTokens int32 `json:"completion_tokens"`
	TotalTokens      int32 `json:"total_tokens"`
}

type Logprobs struct {
	TextOffset    []int32              `json:"text_offset"`
	TokenLogprobs []float32            `json:"token_logprobs"`
	Tokens        []string             `json:"tokens"`
	TopLogprobs   []map[string]float32 `json:"top_logprobs"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}
