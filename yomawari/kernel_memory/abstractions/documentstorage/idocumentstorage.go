package documentstorage

import (
	"context"
	"io"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
)

type IDocumentStorage interface {
	CreateIndexDirectory(ctx context.Context, index string) error
	DeleteIndexDirectory(ctx context.Context, index string) error
	CreateDocumentDirectory(ctx context.Context, index string, documentId string) error
	EmptyDocumentDirectory(ctx context.Context, index string, documentId string) error
	DeleteDocumentDirectory(ctx context.Context, index string, documentId string) error
	WriteFile(ctx context.Context, index string, documentId string, filename string, streamContent io.ReadCloser) error
	ReadFile(ctx context.Context, index string, documentId string, filename string, logErrIfNotFound bool) (*models.StreamableFileContent, error)
}
