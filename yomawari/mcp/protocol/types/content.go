package types

type Content struct {
	Type        string            `json:"type"`
	Text        *string           `json:"text,omitempty"`
	Data        *string           `json:"data,omitempty"`
	MimeType    *string           `json:"mimeType,omitempty"`
	Resource    *ResourceContents `json:"resource,omitempty"`
	Annotations *Annotations      `json:"annotations,omitempty"`
}
