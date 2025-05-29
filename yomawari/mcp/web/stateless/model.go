package stateless

import "github.com/futugyou/yomawari/mcp/protocol"

type StatelessSessionId struct {
	ClientInfo  *protocol.Implementation `json:"clientInfo"`
	UserIdClaim *UserIdClaim             `json:"userIdClaim"`
}

type UserIdClaim struct {
	Issuer string `json:"issuer"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}
