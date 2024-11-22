package openai

import (
	"golang.org/x/exp/slices"
)

const editsPath string = "edits"

var supportedEditModel = []string{
	Text_davinci_edit_001,
	Code_davinci_edit_001,
}

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

type EditService service

func (c *EditService) CreateEdits(request CreateEditsRequest) *CreateEditsResponse {
	result := &CreateEditsResponse{}

	err := validateEditModel(request.Model)
	if err != nil {
		result.Error = err
		return result
	}

	c.client.httpClient.Post(editsPath, request, result)
	return result
}

func validateEditModel(model string) *OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedEditModel, model) {
		return unsupportedTypeError("Model", model, supportedEditModel)
	}

	return nil
}
