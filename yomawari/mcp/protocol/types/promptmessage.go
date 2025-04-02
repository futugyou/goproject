package types

type PromptMessage struct {
	Content Content `json:"content"`
	Role    Role    `json:"role"`
}
