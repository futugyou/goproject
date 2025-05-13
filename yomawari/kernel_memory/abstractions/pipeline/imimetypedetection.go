package pipeline

import (
	"fmt"
	"path/filepath"
)

type IMimeTypeDetection interface {
	GetFileType(filename string) (*string, error)
}

type MimeTypesDetection struct {
}

func (m MimeTypesDetection) GetFileType(filename string) (*string, error) {
	ext := filepath.Ext(filename)
	if condition, ok := FileExtensionToMimeType[ext]; ok {
		return &condition, nil
	}
	return nil, fmt.Errorf("unknown file type: %s", ext)
}

var _ IMimeTypeDetection = &MimeTypesDetection{}
