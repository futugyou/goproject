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
