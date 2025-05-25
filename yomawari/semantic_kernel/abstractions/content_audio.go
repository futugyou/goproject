package abstractions

import "net/url"

type AudioContent struct {
	MimeType     string         `json:"mimeType"`
	ModelId      string         `json:"modelId"`
	Metadata     map[string]any `json:"metadata"`
	Uri          url.URL        `json:"uri"`
	DataUri      string         `json:"-"`
	Data         []byte         `json:"data"`
	InnerContent any            `json:"-"`
}

func (bc *AudioContent) CanRead() bool {
	return bc.Data != nil || bc.DataUri != ""
}

func (AudioContent) Type() string {
	return "audio"
}

func (f AudioContent) ToString() string {
	return string(f.Data)
}

func (c AudioContent) Hash() string {
	return c.ToString()
}

func (f AudioContent) GetInnerContent() any {
	return f.InnerContent
}
