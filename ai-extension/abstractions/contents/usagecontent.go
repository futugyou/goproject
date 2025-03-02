package contents

import (
	"encoding/json"

	"github.com/futugyou/ai-extension/abstractions"
)

// UsageContent represents content that holds usage information.
type UsageContent struct {
	AIContent `json:",inline"`
	Details   abstractions.UsageDetails `json:"details,omitempty"`
}

func (fcc UsageContent) MarshalJSON() ([]byte, error) {
	type Alias UsageContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  "UsageContent",
		Alias: Alias(fcc),
	})
}

func (fcc *UsageContent) UnmarshalJSON(data []byte) error {
	type Alias UsageContent
	aux := &struct {
		Type string `json:"type"`
		Alias
	}{Alias: Alias(*fcc)}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	*fcc = UsageContent(aux.Alias)
	return nil
}
