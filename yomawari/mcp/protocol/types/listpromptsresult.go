package types

import "github.com/futugyou/yomawari/mcp/protocol/messages"

type ListPromptsResult struct {
	messages.PaginatedResult `json:",inline"`
	Prompts                  []Prompt `json:"prompts"`
}
