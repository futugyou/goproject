package types

import "github.com/futugyou/yomawari/mcp/protocol/transport"

type ListPromptsResult struct {
	transport.PaginatedResult `json:",inline"`
	Prompts                   []Prompt `json:"prompts"`
}
