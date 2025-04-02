package types

import "github.com/futugyou/yomawari/mcp/protocol/messages"

type ListResourceTemplatesResult struct {
	messages.PaginatedResult `json:",inline"`
	ResourceTemplates        []ResourceTemplate `json:"resourceTemplates"`
}
