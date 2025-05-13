package contents

import (
	"encoding/json"
)

// FunctionCallContent represents content related to function calls.
type FunctionCallContent struct {
	*AIContent `json:",inline"`
	CallId     string                 `json:"callId,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Arguments  map[string]interface{} `json:"arguments,omitempty"`
	Error      error                  `json:"-"`
}

func CreateFromParsedArguments[TEncoding any](
	encodedArguments TEncoding,
	callId string,
	name string,
	argumentParser func(TEncoding) (map[string]interface{}, error),
) *FunctionCallContent {
	arguments, err := argumentParser(encodedArguments)
	return &FunctionCallContent{
		CallId:    callId,
		Name:      name,
		Arguments: arguments,
		Error:     err,
	}
}

func (ac FunctionCallContent) MarshalJSON() ([]byte, error) {
	type Alias FunctionCallContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "FunctionCallContent",
		Alias: (*Alias)(&ac),
	})
}

func (ac *FunctionCallContent) UnmarshalJSON(data []byte) error {
	type Alias FunctionCallContent
	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(ac),
	}

	return json.Unmarshal(data, aux)
}
