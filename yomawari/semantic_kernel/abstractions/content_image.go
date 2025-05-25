package abstractions

import "net/url"

type ImageContent struct {
	MimeType     string         `json:"mimeType"`
	ModelId      string         `json:"modelId"`
	Metadata     map[string]any `json:"metadata"`
	Uri          url.URL        `json:"uri"`
	DataUri      string         `json:"-"`
	Data         []byte         `json:"data"`
	InnerContent any            `json:"-"`
}

func (bc *ImageContent) CanRead() bool {
	return bc.Data != nil || bc.DataUri != ""
}

func (ImageContent) Type() string {
	return "image"
}

func (f ImageContent) ToString() string {
	return string(f.Data)
}

func (f ImageContent) GetInnerContent() any {
	return f.InnerContent
}

func (c ImageContent) Hash() string {
	return c.ToString()
}
