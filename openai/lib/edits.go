package lib

import (
	"golang.org/x/exp/slices"
)

const editsPath string = "edits"

var supportedEditModel = []string{GPT35_turbo, Code_davinci_edit_001}

type CreateEditsRequest struct {
	Model       string  `json:"model"`
	Input       string  `json:"input,omitempty"`
	Instruction string  `json:"instruction"`
	N           int32   `json:"n,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	Top_p       float32 `json:"top_p,omitempty"`
}

type CreateEditsResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	Object  string       `json:"object,omitempty"`
	Created int          `json:"created,omitempty"`
	Choices []Choices    `json:"choices,omitempty"`
	Usage   *Usage       `json:"usage,omitempty"`
}

func (c *openaiClient) CreateEdits(request CreateEditsRequest) *CreateEditsResponse {
	result := &CreateEditsResponse{}

	err := validateEditModel(request.Model)
	if err != nil {
		result.Error = err
		return result
	}

	c.httpClient.Post(editsPath, request, result)
	return result
}

func validateEditModel(model string) *OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedEditModel, model) {
		return NewError(model, supportedEditModel)
	}

	return nil
}
