package types

import "github.com/futugyou/yomawari/mcp/protocol/transport"

type ListToolsResult struct {
	transport.PaginatedResult `json:",inline"`
	Tools                     []Tool `json:"tools"`
}
