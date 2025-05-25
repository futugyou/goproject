package abstractions

import (
	"encoding/base64"
)

type StreamingFileReferenceContent struct {
	ChoiceIndex  int            `json:"choiceIndex"`
	ModelId      string         `json:"modelId"`
	Metadata     map[string]any `json:"metadata"`
	InnerContent any            `json:"-"`
	FileId       string         `json:"fileId"`
}

func (StreamingFileReferenceContent) Type() string {
	return "streaming-file-reference"
}

func (c StreamingFileReferenceContent) ToString() string {
	return c.FileId
}

func (c StreamingFileReferenceContent) ToByteArray() []byte {
	r, _ := base64.URLEncoding.DecodeString(c.ToString())
	return r
}

func (c StreamingFileReferenceContent) Hash() string {
	return c.FileId
}
