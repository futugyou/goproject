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

var Family = []string{
	GPT35,
	DALLE,
	Whisper,
	Embeddings,
	Codex,
	Moderation,
	GPT3,
}

const GPT35 string = "gpt-3.5"
const DALLE string = "dalle"
const Whisper string = "whisper"
const Embeddings string = "embedding"
const Codex string = "code"
const Moderation string = "moderation"
const GPT3 string = "gpt-3"

var Capability = []string{
	Turbo,
	Davinci,
	Curie,
	Babbage,
	Ada,
}

const Turbo string = "turbo"
const Davinci string = "davinci"
const Curie string = "curie"
const Babbage string = "babbage"
const Ada string = "ada"

const GPT35_turbo string = "gpt-3.5-turbo"
const GPT35_turbo_0301 string = "gpt-3.5-turbo-0301"
const Text_davinci_003 string = "text-davinci-003"
const Text_davinci_002 string = "text-davinci-002"
const Code_davinci_002 string = "code-davinci-002"

const Whisper_1 string = "whisper-1"
const Text_davinci_edit_001 string = "text-davinci-edit-001"
const Code_davinci_edit_001 string = "code-davinci-edit-001"

const Text_moderation_stable string = "text-moderation-stable"
const Text_moderation_latest string = "text-moderation-latest"

var ModelTokenLimitList = map[string]int32{
	GPT35_turbo:      4096,
	GPT35_turbo_0301: 4096,
	Text_davinci_003: 4000,
	Text_davinci_002: 4000,
	Code_davinci_002: 4000,
}
