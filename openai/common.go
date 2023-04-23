package openai

import (
	"math"
)

type Choices struct {
	Text         string                  `json:"text,omitempty"`
	Index        int32                   `json:"index"`
	Logprobs     *Logprobs               `json:"logprobs,omitempty"`
	FinishReason string                  `json:"finish_reason,omitempty"`
	Message      []ChatCompletionMessage `json:"message,omitempty"`
	Delta        ChatCompletionMessage   `json:"delta,omitempty"`
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

type FileModel struct {
	Object        string      `json:"object,omitempty"`
	ID            string      `json:"id,omitempty"`
	Purpose       string      `json:"purpose,omitempty"`
	Filename      string      `json:"filename,omitempty"`
	Bytes         int         `json:"bytes,omitempty"`
	CreatedAt     int         `json:"created_at,omitempty"`
	Status        string      `json:"status,omitempty"`
	StatusDetails interface{} `json:"status_details,omitempty"`
	Owner         string      `json:"owner,omitempty"`
}

var Family = []string{
	GPT4,
	GPT35,
	DALLE,
	Whisper,
	Embeddings,
	Codex,
	Moderation,
	GPT3,
}

const GPT4 string = "gpt-4"
const GPT35 string = "gpt-3.5"
const DALLE string = "dalle"
const Whisper string = "whisper"
const Embeddings string = "embedding"
const Codex string = "code"
const Moderation string = "moderation"
const GPT3 string = "gpt-3"

var GPT4Family = []string{
	GPT_4,
	GPT_4_0314,
	GPT_4_32k,
	GPT_4_32k_0314,
}

const GPT_4 string = "gpt-4"
const GPT_4_0314 string = "gpt-4-0314"
const GPT_4_32k string = "gpt-4-32k"
const GPT_4_32k_0314 string = "gpt-4-32k-0314"

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

var GPT35Family = []string{
	GPT35_turbo,
	GPT35_turbo_0301,
	Text_davinci_003,
	Text_davinci_002,
	Code_davinci_002,
}

const GPT35_turbo string = "gpt-3.5-turbo"
const GPT35_turbo_0301 string = "gpt-3.5-turbo-0301"
const Text_davinci_003 string = "text-davinci-003"
const Text_davinci_002 string = "text-davinci-002"
const Code_davinci_002 string = "code-davinci-002"

var DALLeFamily = []string{
	DALLe_2,
}

const DALLe_2 string = "dalle-2"

var WhisperFamily = []string{
	Whisper_1,
}

const Whisper_1 string = "whisper-1"

var EmbeddingsFamily = []string{
	Text_embedding_ada_002,
	Text_search_ada_doc_001,
}

const Text_embedding_ada_002 string = "text-embedding-ada-002"
const Text_search_ada_doc_001 string = "text-search-ada-doc-001"

var CodexFamily = []string{
	Code_davinci_002,
	Code_cushman_001,
}

const Code_cushman_001 string = "code-cushman-001"

const Text_davinci_edit_001 string = "text-davinci-edit-001"
const Code_davinci_edit_001 string = "code-davinci-edit-001"

var ModerationFamily = []string{
	Text_moderation_stable,
	Text_moderation_latest,
}

const Text_moderation_stable string = "text-moderation-stable"
const Text_moderation_latest string = "text-moderation-latest"

var GPT3Family = []string{
	Text_curie_001,
	Text_babbage_001,
	GPT3_davinci,
	GPT3_curie,
	GPT3_babbage,
	GPT3_ada,
}

const Text_curie_001 string = "text-curie-001"
const Text_babbage_001 string = "text-babbage-001"
const Text_ada_001 string = "text-ada-001"
const GPT3_davinci string = "davinci"
const GPT3_curie string = "curie"
const GPT3_babbage string = "babbage"
const GPT3_ada string = "ada"

var ModelTokenLimitList = map[string]int32{
	GPT_4:            8192,
	GPT_4_0314:       8192,
	GPT_4_32k:        32768,
	GPT_4_32k_0314:   32768,
	GPT35_turbo:      4096,
	GPT35_turbo_0301: 4096,
	Text_davinci_003: 4097,
	Text_davinci_002: 4097,
	Code_davinci_002: 8001,
	Code_cushman_001: 2048,
	Text_curie_001:   2049,
	Text_babbage_001: 2049,
	GPT3_davinci:     2049,
	GPT3_curie:       2049,
	GPT3_babbage:     2049,
	GPT3_ada:         2049,
}

func MaxTokenValidate(model string, token int32, prompt []string) *OpenaiError {
	promptLen := 0
	for _, p := range prompt {
		promptLen += len(p) / 4
	}

	if promptLen > int(token) {
		return &OpenaiError{
			ErrorMessage: "Maximum token limit exceeded",
		}
	}

	if mt, ok := ModelTokenLimitList[model]; ok {
		if token > mt {
			return &OpenaiError{
				ErrorMessage: "Maximum token limit exceeded",
			}
		}
	}

	return nil
}

func CosineSimilarity(listA, listB []float64) float64 {
	count, lenA, lenB := 0, len(listA), len(listB)

	if lenA > lenB {
		count = lenA
	} else {
		count = lenB
	}

	sumA, sumB, sumC := 0.0, 0.0, 0.0
	for i := 0; i < count; i++ {
		if i >= lenA {
			sumC += math.Pow(listB[i], 2)
			continue
		}
		if i >= lenB {
			sumB += math.Pow(listA[i], 2)
			continue
		}
		sumA += listA[i] * listB[i]
		sumB += math.Pow(listA[i], 2)
		sumC += math.Pow(listB[i], 2)
	}

	return sumA / (math.Sqrt(sumB) * math.Sqrt(sumC))
}
