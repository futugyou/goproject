package contents

import "net/url"

type BinaryContent struct {
	MimeType string         `json:"mimeType"`
	ModelId  string         `json:"modelId"`
	Metadata map[string]any `json:"metadata"`
	Uri      url.URL        `json:"uri"`
	DataUri  string         `json:"-"`
	Data     []byte         `json:"data"`
}

func (bc *BinaryContent) CanRead() bool {
	return bc.Data != nil || bc.DataUri != ""
}

func (BinaryContent) Type() string {
	return "binary"
}
