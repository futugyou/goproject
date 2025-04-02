package types

import "github.com/futugyou/yomawari/mcp/protocol/messages"

type RequestParams struct {
	Meta *RequestParamsMetadata `json:"_meta"`
}

type RequestParamsMetadata struct {
	ProgressToken *messages.ProgressToken `json:"progressToken"`
}
