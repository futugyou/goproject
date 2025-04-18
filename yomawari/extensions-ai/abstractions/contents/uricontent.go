package contents

import (
	"encoding/json"
	"strings"
)

type UriContent struct {
	*AIContent `json:",inline"`
	URI        string `json:"url,omitempty"`
	MediaType  string `json:"mediaType,omitempty"`
}

func (dc *UriContent) MediaTypeStartsWith(prefix string) bool {
	return strings.HasPrefix(dc.MediaType, prefix)
}

func (ac UriContent) MarshalJSON() ([]byte, error) {
	type Alias UriContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "UriContent",
		Alias: (*Alias)(&ac),
	})
}

func (ac *UriContent) UnmarshalJSON(data []byte) error {
	type Alias UriContent
	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(ac),
	}

	return json.Unmarshal(data, aux)
}
