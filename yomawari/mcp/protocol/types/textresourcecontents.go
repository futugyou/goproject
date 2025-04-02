package types

import "encoding/json"

type TextResourceContents struct {
	BaseResourceContents `json:",inline"`
	Text                 string `json:"text"`
}

func (b TextResourceContents) MarshalJSON() ([]byte, error) {
	type Alias TextResourceContents
	return json.Marshal(&struct {
		Alias
		Type string `json:"type"`
	}{
		Alias: Alias(b),
		Type:  "text",
	})
}

func (b *TextResourceContents) UnmarshalJSON(data []byte) error {
	type Alias TextResourceContents
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	return nil
}
