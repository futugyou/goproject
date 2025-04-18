package contents

import "encoding/json"

// FunctionResultContent represents the result of a function call.
type FunctionResultContent struct {
	*AIContent `json:",inline"`
	CallId     string      `json:"callId,omitempty"`
	Result     interface{} `json:"result,omitempty"`
	Error      error       `json:"-"`
}

func (ac FunctionResultContent) MarshalJSON() ([]byte, error) {
	type Alias FunctionResultContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "FunctionResultContent",
		Alias: (*Alias)(&ac),
	})
}

func (ac *FunctionResultContent) UnmarshalJSON(data []byte) error {
	type Alias FunctionResultContent
	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(ac),
	}

	return json.Unmarshal(data, aux)
}
