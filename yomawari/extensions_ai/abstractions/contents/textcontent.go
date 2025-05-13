package contents

import "encoding/json"

// TextContent represents text-based content.
type TextContent struct {
	*AIContent `json:",inline"`
	Text       string `json:"text,omitempty"`
}

func NewTextContent(text string) *TextContent {
	return &TextContent{
		AIContent: NewAIContent(nil, nil),
		Text:      text,
	}
}

func NewTextContentWithRefusal(text string, refusal string) *TextContent {
	c := NewTextContent(text)
	if len(refusal) > 0 {
		c.AdditionalProperties["refusal"] = refusal
	}
	return c
}

func (ac TextContent) MarshalJSON() ([]byte, error) {
	type Alias TextContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "TextContent",
		Alias: (*Alias)(&ac),
	})
}

func (ac *TextContent) UnmarshalJSON(data []byte) error {
	type Alias TextContent
	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(ac),
	}

	return json.Unmarshal(data, aux)
}
