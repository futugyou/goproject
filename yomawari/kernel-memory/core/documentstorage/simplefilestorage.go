package documentstorage

import (
	"context"
	"io"

	"github.com/futugyou/yomawari/kernel-memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel-memory/abstractions/models"
)

type SimpleFileStorage struct {
	// TODO: implement
}

// CreateDocumentDirectory implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) CreateDocumentDirectory(ctx context.Context, index string, documentId string) error {
	panic("unimplemented")
}

// CreateIndexDirectory implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) CreateIndexDirectory(ctx context.Context, index string) error {
	panic("unimplemented")
}

// DeleteDocumentDirectory implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) DeleteDocumentDirectory(ctx context.Context, index string, documentId string) error {
	panic("unimplemented")
}

// DeleteIndexDirectory implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) DeleteIndexDirectory(ctx context.Context, index string) error {
	panic("unimplemented")
}

// EmptyDocumentDirectory implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) EmptyDocumentDirectory(ctx context.Context, index string, documentId string) error {
	panic("unimplemented")
}

// ReadFile implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) ReadFile(ctx context.Context, index string, documentId string, filename string, logErrIfNotFound bool) (*models.StreamableFileContent, error) {
	panic("unimplemented")
}

// WriteFile implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) WriteFile(ctx context.Context, index string, documentId string, filename string, streamContent io.ReadCloser) error {
	panic("unimplemented")
}

var _ documentstorage.IDocumentStorage = (*SimpleFileStorage)(nil)
