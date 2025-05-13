package filesystem

import (
	"context"
	"io"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
)

type IFileSystem interface {
	// Volume API
	CreateVolume(ctx context.Context, volume string) error
	VolumeExists(ctx context.Context, volume string) bool
	DeleteVolume(ctx context.Context, volume string) error
	ListVolumes(ctx context.Context) ([]string, error)

	// Directory API
	CreateDirectory(ctx context.Context, volume, relPath string) error
	DeleteDirectory(ctx context.Context, volume, relPath string) error

	// File API
	WriteFile(ctx context.Context, volume, relPath, fileName string, streamContent io.Reader) error
	WriteFileAsText(ctx context.Context, volume, relPath, fileName, data string) error
	FileExists(ctx context.Context, volume, relPath, fileName string) bool

	ReadFileAsBinary(ctx context.Context, volume, relPath, fileName string) ([]byte, error)
	ReadFileInfo(ctx context.Context, volume, relPath, fileName string) (*models.StreamableFileContent, error)
	ReadFileAsText(ctx context.Context, volume, relPath, fileName string) (*string, error)
	ReadAllFilesAsText(ctx context.Context, volume, relPath string) (map[string]string, error)
	GetAllFileNames(ctx context.Context, volume, relPath string) ([]string, error)

	DeleteFile(ctx context.Context, volume, relPath, fileName string) error
}
