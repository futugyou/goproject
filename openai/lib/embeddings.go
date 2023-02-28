package lib

const embeddingsPath string = "embeddings"

type CreateEmbeddingsRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
	User  float32  `json:"user,omitempty"`
}

type CreateEmbeddingsResponse struct {
	Error  *OpenaiError     `json:"error,omitempty"`
	Object string           `json:"object,omitempty"`
	Data   []embeddingsData `json:"data,omitempty"`
	Model  string           `json:"model,omitempty"`
	Usage  *usage           `json:"usage,omitempty"`
}
type embeddingsData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

func (client *openaiClient) CreateEmbeddings(request CreateEmbeddingsRequest) *CreateEmbeddingsResponse {
	result := &CreateEmbeddingsResponse{}
	client.Post(embeddingsPath, request, result)
	return result
}
