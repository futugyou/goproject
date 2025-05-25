package abstractions

import (
	"encoding/base64"
	"fmt"
)

type StreamingAnnotationContent struct {
	ChoiceIndex  int            `json:"choiceIndex"`
	ModelId      string         `json:"modelId"`
	Metadata     map[string]any `json:"metadata"`
	InnerContent any            `json:"-"`
	FileId       string         `json:"fileId"`
	Quote        string         `json:"quote"`
	StartIndex   int            `json:"startIndex"`
	EndIndex     int            `json:"endIndex"`
}

func (StreamingAnnotationContent) Type() string {
	return "streaming-annotation"
}

func (c StreamingAnnotationContent) ToString() string {
	hasFileId := len(c.FileId) > 0

	if hasFileId {
		return fmt.Sprintf("%s: %s", c.Quote, c.FileId)
	}

	return c.Quote
}

func (c StreamingAnnotationContent) ToByteArray() []byte {
	r, _ := base64.URLEncoding.DecodeString(c.ToString())
	return r
}

func (c StreamingAnnotationContent) Hash() string {
	return c.ToString()
}
