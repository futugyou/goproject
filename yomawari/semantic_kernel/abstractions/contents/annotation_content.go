package contents

type AnnotationContent struct {
	FileId       string         `json:"fileId"`
	Quote        string         `json:"quote"`
	StartIndex   int            `json:"startIndex"`
	EndIndex     int            `json:"endIndex"`
	MimeType     string         `json:"mimeType"`
	ModelId      string         `json:"modelId"`
	Metadata     map[string]any `json:"metadata"`
	InnerContent any            `json:"-"`
}

func (AnnotationContent) Type() string {
	return "annotation"
}

func (f AnnotationContent) ToString() string {
	return f.Quote
}

func (f AnnotationContent) GetInnerContent() any {
	return f.InnerContent
}
