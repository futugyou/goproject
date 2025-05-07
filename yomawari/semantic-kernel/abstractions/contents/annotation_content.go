package contents

type AnnotationContent struct {
	FileId     string         `json:"fileId"`
	Quote      string         `json:"quote"`
	StartIndex int            `json:"startIndex"`
	EndIndex   int            `json:"endIndex"`
	MimeType   string         `json:"mimeType"`
	ModelId    string         `json:"modelId"`
	Metadata   map[string]any `json:"metadata"`
}

func (AnnotationContent) Type() string {
	return "annotation"
}
