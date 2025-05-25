package abstractions

import "encoding/base64"

type TextContent struct {
	MimeType     string           `json:"mimeType"`
	ModelId      string           `json:"modelId"`
	Metadata     map[string]any   `json:"metadata"`
	InnerContent any              `json:"-"`
	Text         string           `json:"text"`
	Encoding     *base64.Encoding `json:"-"`
}

func (TextContent) Type() string {
	return "text"
}

func (f TextContent) ToString() string {
	return f.Text
}

func (f TextContent) GetInnerContent() any {
	return f.InnerContent
}

func (c TextContent) Hash() string {
	return c.Text
}
