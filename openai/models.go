package openai

import (
	"fmt"
)

const listModelsPath string = "models"
const retrieveModelPath string = "models/%s"

type ListModelResponse struct {
	Error  *OpenaiError `json:"error,omitempty"`
	Object string       `json:"object,omitempty"`
	Datas  []model      `json:"data,omitempty"`
}

type ModelResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	model
}

type model struct {
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

func (c *openaiClient) ListModels() *ListModelResponse {
	result := &ListModelResponse{}
	c.httpClient.Get(listModelsPath, result)
	return result
}

func (c *openaiClient) RetrieveModel(model string) *ModelResponse {
	result := &ModelResponse{}
	c.httpClient.Get(fmt.Sprintf(retrieveModelPath, model), result)
	return result
}
