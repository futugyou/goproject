package protocol

var resourceFactories = map[string]func() IResourceContents{
	"blob": func() IResourceContents { return &BlobResourceContents{} },
	"text": func() IResourceContents { return &TextResourceContents{} },
}

type IResourceContents interface {
	GetUri() string
	GetMimeType() *string
}

type BaseResourceContents struct {
	Uri      string  `json:"uri"`
	MimeType *string `json:"mimeType,omitempty"`
}

// IsResourceContents implements IResourceContents.
func (r *BaseResourceContents) GetMimeType() *string {
	return r.MimeType
}

// IsResourceContents implements IResourceContents.
func (r *BaseResourceContents) GetUri() string {
	return r.Uri
}
