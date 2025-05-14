package stateless

import "github.com/futugyou/yomawari/mcp/protocol/types"

type StatelessSessionId struct {
	ClientInfo  *types.Implementation `json:"clientInfo"`
	UserIdClaim *UserIdClaim          `json:"userIdClaim"`
}

type UserIdClaim struct {
	Issuer string `json:"issuer"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}
