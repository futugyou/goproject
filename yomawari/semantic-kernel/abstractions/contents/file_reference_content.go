package contents

type FileReferenceContent struct {
	MimeType string         `json:"mimeType"`
	ModelId  string         `json:"modelId"`
	Metadata map[string]any `json:"metadata"`
	FileId   string         `json:"fileId"`
	Tools    []string       `json:"tools"`
}

func (FileReferenceContent) Type() string {
	return "fileReference"
}
