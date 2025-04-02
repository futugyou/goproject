package types

type ListPromptsResult struct {
	PaginatedResult
	Prompts []Prompt `json:"prompts"`
}
