package contents

import "net/url"

type AudioContent struct {
	MimeType string         `json:"mimeType"`
	ModelId  string         `json:"modelId"`
	Metadata map[string]any `json:"metadata"`
	Uri      url.URL        `json:"uri"`
	DataUri  string         `json:"-"`
	Data     []byte         `json:"data"`
}

func (bc *AudioContent) CanRead() bool {
	return bc.Data != nil || bc.DataUri != ""
}

func (AudioContent) Type() string {
	return "audio"
}
