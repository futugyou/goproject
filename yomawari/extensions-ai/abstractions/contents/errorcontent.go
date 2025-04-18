package contents

import "encoding/json"

type ErrorContent struct {
	*AIContent `json:",inline"`
	Message    string  `json:"message"`
	ErrorCode  *string `json:"errorCode,omitempty"`
	Details    *string `json:"details,omitempty"`
}

func (ac ErrorContent) MarshalJSON() ([]byte, error) {
	type Alias ErrorContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "ErrorContent",
		Alias: (*Alias)(&ac),
	})
}

func (ac *ErrorContent) UnmarshalJSON(data []byte) error {
	type Alias ErrorContent
	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(ac),
	}

	return json.Unmarshal(data, aux)
}
