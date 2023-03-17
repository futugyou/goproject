package lib

import (
	e "openai/lib/internal"

	"golang.org/x/exp/slices"
)

const embeddingsPath string = "embeddings"

var supportedEmbeddingsModel = []string{
	Text_embedding_ada_002,
	Text_search_ada_doc_001,
}

type CreateEmbeddingsRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
	User  float32  `json:"user,omitempty"`
}

type CreateEmbeddingsResponse struct {
	Error  *e.OpenaiError   `json:"error,omitempty"`
	Object string           `json:"object,omitempty"`
	Data   []EmbeddingsData `json:"data,omitempty"`
	Model  string           `json:"model,omitempty"`
	Usage  *Usage           `json:"usage,omitempty"`
}

type EmbeddingsData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

func (c *openaiClient) CreateEmbeddings(request CreateEmbeddingsRequest) *CreateEmbeddingsResponse {
	result := &CreateEmbeddingsResponse{}
	err := validateEmbeddingsModel(request.Model)
	if err != nil {
		result.Error = err
		return result
	}

	c.httpClient.Post(embeddingsPath, request, result)
	return result
}

func validateEmbeddingsModel(model string) *e.OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedEmbeddingsModel, model) {
		return e.UnsupportedTypeError("Model", model, supportedEmbeddingsModel)
	}

	return nil
}
