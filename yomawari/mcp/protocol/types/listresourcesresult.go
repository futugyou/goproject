package types

import "github.com/futugyou/yomawari/mcp/protocol/messages"

type ListResourcesResult struct {
	messages.PaginatedResult `json:",inline"`
	Resources                []Resource `json:"resources"`
}
