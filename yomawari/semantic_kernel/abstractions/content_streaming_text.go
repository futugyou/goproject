package abstractions

import "encoding/base64"

type StreamingTextContent struct {
	ChoiceIndex  int              `json:"choiceIndex"`
	ModelId      string           `json:"modelId"`
	Metadata     map[string]any   `json:"metadata"`
	InnerContent any              `json:"-"`
	Text         string           `json:"text"`
	Encoding     *base64.Encoding `json:"-"`
}

func (StreamingTextContent) Type() string {
	return "streaming-function-call-update"
}

func (c StreamingTextContent) ToString() string {
	return c.Text
}

func (c StreamingTextContent) ToByteArray() []byte {
	encoding := base64.URLEncoding
	if c.Encoding != nil {
		encoding = c.Encoding
	}
	r, _ := encoding.DecodeString(c.ToString())
	return r
}

func (c StreamingTextContent) Hash() string {
	return c.Text
}
