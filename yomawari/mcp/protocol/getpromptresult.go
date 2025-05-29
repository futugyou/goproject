package protocol

type GetPromptResult struct {
	Description *string         `json:"description"`
	Messages    []PromptMessage `json:"messages"`
}
