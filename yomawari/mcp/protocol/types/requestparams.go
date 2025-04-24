package types

import "github.com/futugyou/yomawari/mcp/protocol/transport"

type RequestParams struct {
	Meta *RequestParamsMetadata `json:"_meta"`
}

type RequestParamsMetadata struct {
	ProgressToken *transport.ProgressToken `json:"progressToken"`
}
