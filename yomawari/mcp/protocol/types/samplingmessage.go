package types

type SamplingMessage struct {
	Content Content `json:"content"`
	Role    Role    `json:"role"`
}
