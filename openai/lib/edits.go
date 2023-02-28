package lib

const editsPath string = "edits"

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
	Choices []choices    `json:"choices,omitempty"`
	Usage   *usage       `json:"usage,omitempty"`
}

func (client *openaiClient) CreateEdits(request CreateEditsRequest) *CreateEditsResponse {
	result := &CreateEditsResponse{}
	client.Post(editsPath, request, result)
	return result
}
