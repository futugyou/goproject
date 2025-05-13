package ollama

type OllamaEmbeddingRequest struct {
	Model     string
	Input     []string
	Options   *OllamaRequestOptions
	Truncate  bool
	KeepAlive *int64
}
