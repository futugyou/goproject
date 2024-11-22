package openai

import (
	"context"
	"fmt"
)

const listModelsPath string = "models"
const retrieveModelPath string = "models/%s"

type ListModelResponse struct {
	Error  *OpenaiError `json:"error,omitempty"`
	Object string       `json:"object,omitempty"`
	Datas  []Model      `json:"data,omitempty"`
}

type ModelResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	Model
}

type Model struct {
	Error      *OpenaiError `json:"error,omitempty"`
	ID         string       `json:"id"`
	Object     string       `json:"object"`
	Created    int32        `json:"created"`
	OwnedBy    string       `json:"owned_by"`
	Permission []permission `json:"permission"`
	Root       string       `json:"root"`
	Parent     interface{}  `json:"parent"`
}

type permission struct {
	ID                 string      `json:"id"`
	Object             string      `json:"object"`
	Created            int32       `json:"created"`
	AllowCreateEngine  bool        `json:"allow_create_engine"`
	AllowSampling      bool        `json:"allow_sampling"`
	AllowLogprobs      bool        `json:"allow_logprobs"`
	AllowSearchIndices bool        `json:"allow_search_indices"`
	AllowView          bool        `json:"allow_view"`
	AllowFineTuning    bool        `json:"allow_fine_tuning"`
	Organization       string      `json:"organization"`
	Group              interface{} `json:"group"`
	IsBlocking         bool        `json:"is_blocking"`
}

type ModelService service

func (c *ModelService) ListModels(ctx context.Context) *ListModelResponse {
	result := &ListModelResponse{}
	if err := c.client.httpClient.Get(ctx, listModelsPath, result); err != nil {
		result.Error = err
	}
	return result
}

func (c *ModelService) RetrieveModel(ctx context.Context, model string) *ModelResponse {
	result := &ModelResponse{}
	if err := c.client.httpClient.Get(ctx, fmt.Sprintf(retrieveModelPath, model), result); err != nil {
		result.Error = err
	}
	return result
}
