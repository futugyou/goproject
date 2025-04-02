package types

type BlobResourceContents struct {
	ResourceContents
	Blob string `json:"blob"`
}
