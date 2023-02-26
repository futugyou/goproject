package pkg

const editsPath string = "completions"

type CreateEditsRequest struct {
	Model       string  `json:"model"`
	Input       string  `json:"input,omitempty"`
	Instruction string  `json:"instruction"`
	N           int32   `json:"n,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	Top_p       float32 `json:"top_p,omitempty"`
}
