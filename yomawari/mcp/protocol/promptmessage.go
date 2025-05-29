package protocol

type PromptMessage struct {
	Content Content `json:"content"`
	Role    Role    `json:"role"`
}
