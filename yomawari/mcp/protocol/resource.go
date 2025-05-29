package protocol

type Resource struct {
	Uri         string       `json:"uri"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	MimeType    *string      `json:"mimeType,omitempty"`
	Size        *float32     `json:"size,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}
