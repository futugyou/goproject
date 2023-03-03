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

const GPT35_turbo string = "gpt-3.5-turbo"
const GPT35_turbo_0301 string = "gpt-3.5-turbo-0301"
const Whisper_1 string = "whisper-1"
const Text_davinci_edit_001 string = "text-davinci-edit-001"
const Code_davinci_edit_001 string = "code-davinci-edit-001"

const Ada string = "ada"
const Babbage string = "babbage"
const Curie string = "curie"
const Davinci string = "davinci"
