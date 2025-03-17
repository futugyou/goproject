package contents

import "encoding/json"

// TextContent represents text-based content.
type TextContent struct {
	AIContent `json:",inline"`
	Text      string `json:"text,omitempty"`
}

func NewTextContent(text string) *TextContent {
	return &TextContent{
		AIContent: AIContent{AdditionalProperties: make(map[string]interface{})},
		Text:      text,
	}
}

func (fcc TextContent) MarshalJSON() ([]byte, error) {
	type Alias TextContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  "TextContent",
		Alias: Alias(fcc),
	})
}

func (fcc *TextContent) UnmarshalJSON(data []byte) error {
	type Alias TextContent
	aux := &struct {
		Type string `json:"type"`
		Alias
	}{Alias: Alias(*fcc)}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	*fcc = TextContent(aux.Alias)
	return nil
}
