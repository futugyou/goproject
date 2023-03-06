package models

type QuestionAnswer struct {
	Prompt           string   `json:"prompt"`
	MaxTokens        int32    `json:"max_tokens"`
	Temperature      float32  `json:"temperature"`
	TopP             float32  `json:"top_p"`
	FrequencyPenalty float32  `json:"frequency_penalty"`
	PresencePenalty  float32  `json:"presence_penalty"`
	BestOf           int      `json:"best_of" valid:"Range(1, 20)"`
	Echo             bool     `json:"echo"`
	Logprobs         int      `json:"logprobs"`
	Stop             []string `json:"stop"`
	Stream           bool     `json:"stream"`
}
