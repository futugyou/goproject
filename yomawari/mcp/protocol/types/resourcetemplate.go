package types

type ResourceTemplate struct {
	UriTemplate string       `json:"uriTemplate"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	MimeType    *string      `json:"mimeType"`
	Annotations *Annotations `json:"annotations"`
}
