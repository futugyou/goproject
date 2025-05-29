package protocol

import "encoding/json"

type TextResourceContents struct {
	BaseResourceContents `json:",inline"`
	Text                 string `json:"text"`
}

func (b TextResourceContents) MarshalJSON() ([]byte, error) {
	type Alias TextResourceContents
	return json.Marshal(&struct {
		Alias
		ResourceType string `json:"resource_type"`
	}{
		Alias:        Alias(b),
		ResourceType: "text",
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
