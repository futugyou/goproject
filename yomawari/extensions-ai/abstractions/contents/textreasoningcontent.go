package contents

import "encoding/json"

type TextReasoningContent struct {
	*AIContent `json:",inline"`
	Text       string `json:"text"`
}

func NewTextReasoningContent(text string) *TextReasoningContent {
	return &TextReasoningContent{
		AIContent: NewAIContent(nil, nil),
		Text:      text,
	}
}

func (ac TextReasoningContent) MarshalJSON() ([]byte, error) {
	type Alias TextReasoningContent
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "TextReasoningContent",
		Alias: (*Alias)(&ac),
	})
}

func (ac *TextReasoningContent) UnmarshalJSON(data []byte) error {
	type Alias TextReasoningContent
	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(ac),
	}

	return json.Unmarshal(data, aux)
}
