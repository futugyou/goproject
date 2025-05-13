package contents

import "net/url"

// A general binary data processing method is needed for use in audio/image, etc.
type BinaryContent struct {
	MimeType     string         `json:"mimeType"`
	ModelId      string         `json:"modelId"`
	Metadata     map[string]any `json:"metadata"`
	Uri          url.URL        `json:"uri"`
	DataUri      string         `json:"-"`
	Data         []byte         `json:"data"`
	InnerContent any            `json:"-"`
}

func (bc *BinaryContent) CanRead() bool {
	return bc.Data != nil || bc.DataUri != ""
}

func (BinaryContent) Type() string {
	return "binary"
}

func (f BinaryContent) ToString() string {
	return string(f.Data)
}

func (f BinaryContent) GetInnerContent() any {
	return f.InnerContent
}
