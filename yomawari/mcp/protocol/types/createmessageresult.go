package types

type CreateMessageResult struct {
	Content    Content `json:"content"`
	Model      string  `json:"model"`
	StopReason *string `json:"stopReason"`
	Role       string  `json:"role"`
}
