package ollama

type OllamaRequestOptions struct {
	EmbeddingOnly    *bool    `json:"embedding_only,omitempty"`
	F16KV            *bool    `json:"f16_kv,omitempty"`
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	LogitsAll        *bool    `json:"logits_all,omitempty"`
	LowVRAM          *bool    `json:"low_vram,omitempty"`
	MainGPU          *int     `json:"main_gpu,omitempty"`
	MinP             *float64 `json:"min_p,omitempty"`
	Mirostat         *int     `json:"mirostat,omitempty"`
	MirostatEta      *float64 `json:"mirostat_eta,omitempty"`
	MirostatTau      *float64 `json:"mirostat_tau,omitempty"`
	NumBatch         *int     `json:"num_batch,omitempty"`
	NumCtx           *int     `json:"num_ctx,omitempty"`
	NumGPU           *int     `json:"num_gpu,omitempty"`
	NumKeep          *int     `json:"num_keep,omitempty"`
	NumPredict       *int64   `json:"num_predict,omitempty"`
	NumThread        *int     `json:"num_thread,omitempty"`
	NUMA             *bool    `json:"numa,omitempty"`
	PenalizeNewline  *bool    `json:"penalize_newline,omitempty"`
	PresencePenalty  *float64 `json:"presence_penalty,omitempty"`
	RepeatLastN      *int     `json:"repeat_last_n,omitempty"`
	RepeatPenalty    *float64 `json:"repeat_penalty,omitempty"`
	Seed             *int64   `json:"seed,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	Temperature      *float64 `json:"temperature,omitempty"`
	TFSZ             *float64 `json:"tfs_z,omitempty"`
	TopK             *int     `json:"top_k,omitempty"`
	TopP             *float64 `json:"top_p,omitempty"`
	TypicalP         *float64 `json:"typical_p,omitempty"`
	UseMLock         *bool    `json:"use_mlock,omitempty"`
	UseMMap          *bool    `json:"use_mmap,omitempty"`
	VocabOnly        *bool    `json:"vocab_only,omitempty"`
}
