package contents

import (
	"encoding/json"

	"github.com/futugyou/yomawari/extensions_ai/abstractions"
)

// UsageContent represents content that holds usage information.
type UsageContent struct {
	*AIContent `json:",inline"`
	Details    abstractions.UsageDetails `json:"details,omitempty"`
}

func NewUsageContent(details abstractions.UsageDetails) *UsageContent {
	return &UsageContent{
		AIContent: NewAIContent(nil, nil),
		Details:   details,
	}
}

func (ac UsageContent) MarshalJSON() ([]byte, error) {
	type Alias UsageContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "UsageContent",
		Alias: (*Alias)(&ac),
	})
}

func (ac *UsageContent) UnmarshalJSON(data []byte) error {
	type Alias UsageContent
	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(ac),
	}

	return json.Unmarshal(data, aux)
}
