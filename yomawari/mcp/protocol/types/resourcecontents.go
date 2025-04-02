package types

type ResourceContents struct {
	Uri      string  `json:"uri"`
	MimeType *string `json:"mimeType"`
}

// TODO: Marshal/Unmarshal
