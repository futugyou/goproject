package abstractions

import (
	"fmt"
	"net/url"
	"os"
)

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

func (c BinaryContent) Hash() string {
	return c.ToString()
}

func (f BinaryContent) GetInnerContent() any {
	return f.InnerContent
}

func (f *BinaryContent) WriteToFile(filePath string, overwrite bool) error {
	if len(filePath) == 0 {
		return fmt.Errorf("filePath is empty")
	}

	_, err := os.Stat(filePath)
	if err == nil && !overwrite {
		return fmt.Errorf("file exists")
	}

	if !f.CanRead() {
		return fmt.Errorf("no content to write to file")
	}

	return os.WriteFile(filePath, f.Data, 0644)
}
