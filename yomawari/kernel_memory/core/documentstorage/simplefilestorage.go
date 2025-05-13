package documentstorage

import (
	"context"
	"io"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/constant"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/core/filesystem"
)

type SimpleFileStorage struct {
	fileSystem filesystem.IFileSystem
}

func NewSimpleFileStorage(config SimpleFileStorageConfig, mimeTypeDetection pipeline.IMimeTypeDetection) *SimpleFileStorage {
	var fileSystem filesystem.IFileSystem
	if config.StorageType == filesystem.Volatile {
		fileSystem = filesystem.NewVolatileFileSystem(mimeTypeDetection)
	} else {
		fileSystem = filesystem.NewDiskFileSystem(config.Directory, mimeTypeDetection)
	}
	return &SimpleFileStorage{
		fileSystem: fileSystem,
	}
}

// CreateDocumentDirectory implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) CreateDocumentDirectory(ctx context.Context, index string, documentId string) error {
	return s.fileSystem.CreateDirectory(ctx, index, documentId)
}

// CreateIndexDirectory implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) CreateIndexDirectory(ctx context.Context, index string) error {
	return s.fileSystem.CreateVolume(ctx, index)
}

// DeleteDocumentDirectory implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) DeleteDocumentDirectory(ctx context.Context, index string, documentId string) error {
	return s.fileSystem.DeleteDirectory(ctx, index, documentId)
}

// DeleteIndexDirectory implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) DeleteIndexDirectory(ctx context.Context, index string) error {
	return s.fileSystem.DeleteVolume(ctx, index)
}

// EmptyDocumentDirectory implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) EmptyDocumentDirectory(ctx context.Context, index string, documentId string) error {
	files, err := s.fileSystem.GetAllFileNames(ctx, index, documentId)
	if err != nil {
		return err
	}
	for _, fileName := range files {
		if fileName == constant.PipelineStatusFilename {
			continue
		}
		s.fileSystem.DeleteFile(ctx, index, documentId, fileName)
	}
	return nil
}

// ReadFile implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) ReadFile(ctx context.Context, index string, documentId string, filename string, logErrIfNotFound bool) (*models.StreamableFileContent, error) {
	return s.fileSystem.ReadFileInfo(ctx, index, documentId, filename)
}

// WriteFile implements documentstorage.IDocumentStorage.
func (s *SimpleFileStorage) WriteFile(ctx context.Context, index string, documentId string, filename string, streamContent io.ReadCloser) error {
	err := s.fileSystem.CreateDirectory(ctx, index, documentId)
	if err != nil {
		return err
	}
	return s.fileSystem.WriteFile(ctx, index, documentId, filename, streamContent)
}

var _ documentstorage.IDocumentStorage = (*SimpleFileStorage)(nil)
