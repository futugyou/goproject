package protocol

type Prompt struct {
	Arguments   []PromptArgument `json:"arguments"`
	Description *string          `json:"description"`
	Name        string           `json:"name"`
}
