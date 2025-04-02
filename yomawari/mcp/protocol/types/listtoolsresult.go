package types

import "github.com/futugyou/yomawari/mcp/protocol/messages"

type ListToolsResult struct {
	messages.PaginatedResult `json:",inline"`
	Tools                    []Tool `json:"tools"`
}
