package abstractions

type FileReferenceContent struct {
	MimeType     string         `json:"mimeType"`
	ModelId      string         `json:"modelId"`
	Metadata     map[string]any `json:"metadata"`
	FileId       string         `json:"fileId"`
	Tools        []string       `json:"tools"`
	InnerContent any            `json:"-"`
}

func (FileReferenceContent) Type() string {
	return "fileReference"
}

func (f FileReferenceContent) ToString() string {
	return f.FileId
}

func (f FileReferenceContent) GetInnerContent() any {
	return f.InnerContent
}

func (c FileReferenceContent) Hash() string {
	return c.ToString()
}
