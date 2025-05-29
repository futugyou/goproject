package protocol

type PromptArgument struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Required    *bool   `json:"required"`
}
