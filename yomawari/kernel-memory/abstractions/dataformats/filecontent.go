package dataformats

type FileContent struct {
	Sections []Chunk `json:"sections"`
	MimeType string  `json:"mimetype"`
}
