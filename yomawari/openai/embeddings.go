package openai

import (
	"context"

	"golang.org/x/exp/slices"

	"github.com/futugyou/yomawari/common/errorutils"
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
	Error  *errorutils.OpenaiError `json:"error,omitempty"`
	Object string                  `json:"object,omitempty"`
	Data   []EmbeddingsData        `json:"data,omitempty"`
	Model  string                  `json:"model,omitempty"`
	Usage  *Usage                  `json:"usage,omitempty"`
}

type EmbeddingsData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

type EmbeddingService service

func (c *EmbeddingService) CreateEmbeddings(ctx context.Context, request CreateEmbeddingsRequest) *CreateEmbeddingsResponse {
	result := &CreateEmbeddingsResponse{}
	err := validateEmbeddingsModel(request.Model)
	if err != nil {
		result.Error = err
		return result
	}

	c.client.httpClient.Post(ctx, embeddingsPath, request, result)
	return result
}

func validateEmbeddingsModel(model string) *errorutils.OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedEmbeddingsModel, model) {
		return errorutils.UnsupportedTypeError("Model", model, supportedEmbeddingsModel)
	}

	return nil
}
