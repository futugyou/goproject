package types

import "github.com/futugyou/yomawari/mcp/protocol/transport"

type ListResourcesResult struct {
	transport.PaginatedResult `json:",inline"`
	Resources                 []Resource `json:"resources"`
}
