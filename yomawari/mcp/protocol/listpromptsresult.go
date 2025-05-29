package protocol

type ListPromptsResult struct {
	PaginatedResult `json:",inline"`
	Prompts         []Prompt `json:"prompts"`
}
