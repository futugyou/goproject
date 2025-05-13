package models

import "io"

type DocumentUploadRequest struct {
	Index      string
	DocumentId string
	Tags       *TagCollection
	Files      []UploadedFile
	Steps      []string
}

func NewDocumentUploadRequest(document Document, index *string, steps []string) *DocumentUploadRequest {
	d := &DocumentUploadRequest{
		DocumentId: document.Id,
		Tags:       document.Tags,
		Files:      []UploadedFile{},
		Steps:      steps,
	}

	if index != nil {
		d.Index = *index
	}

	if document.Files != nil {
		for _, v := range document.Files.GetStreams() {
			updatedFile := UploadedFile{
				FileName:    v.Name,
				FileContent: v.Content,
			}
			d.Files = append(d.Files, updatedFile)
		}
	}

	return d
}

type UploadedFile struct {
	FileName    string
	FileContent io.ReadCloser
}
