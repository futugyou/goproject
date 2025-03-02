package contents

import "encoding/json"

// FunctionResultContent represents the result of a function call.
type FunctionResultContent struct {
	AIContent `json:",inline"`
	CallId    string      `json:"callId,omitempty"`
	Result    interface{} `json:"result,omitempty"`
	Error     error       `json:"-"`
}

func (fcc FunctionResultContent) MarshalJSON() ([]byte, error) {
	type Alias FunctionResultContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  "FunctionResultContent",
		Alias: Alias(fcc),
	})
}

func (fcc *FunctionResultContent) UnmarshalJSON(data []byte) error {
	type Alias FunctionResultContent
	aux := &struct {
		Type string `json:"type"`
		Alias
	}{Alias: Alias(*fcc)}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	*fcc = FunctionResultContent(aux.Alias)
	return nil
}
