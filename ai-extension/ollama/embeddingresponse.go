package ollama

type OllamaEmbeddingResponse struct {
	Model           *string     `json:"model,omitempty"`
	Embeddings      [][]float64 `json:"embeddings,omitempty"`
	TotalDuration   *int64      `json:"total_duration,omitempty"`
	LoadDuration    *int64      `json:"load_duration,omitempty"`
	PromptEvalCount *int64      `json:"prompt_eval_count,omitempty"`
	Error           *string     `json:"-"`
}
