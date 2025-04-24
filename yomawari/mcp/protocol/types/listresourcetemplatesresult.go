package types

import "github.com/futugyou/yomawari/mcp/protocol/transport"

type ListResourceTemplatesResult struct {
	transport.PaginatedResult `json:",inline"`
	ResourceTemplates         []ResourceTemplate `json:"resourceTemplates"`
}
