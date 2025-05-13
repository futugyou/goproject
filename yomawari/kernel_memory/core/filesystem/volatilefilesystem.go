package filesystem

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
)

var (
	sf_singletons sync.Map
)

type VolatileFileSystem struct {
	mutex             sync.Mutex
	mimeTypeDetection pipeline.IMimeTypeDetection
	_volumes          map[string]map[string][]byte
}

func NewVolatileFileSystem(mimeTypeDetection pipeline.IMimeTypeDetection) *VolatileFileSystem {
	if mimeTypeDetection == nil {
		mimeTypeDetection = &pipeline.MimeTypesDetection{}
	}
	return &VolatileFileSystem{
		mimeTypeDetection: mimeTypeDetection,
	}
}

func GetVolatileFileSystemInstance(directory string, mimeTypeDetection pipeline.IMimeTypeDetection) *VolatileFileSystem {
	directory = strings.Trim(directory, "/\\")
	directory = strings.ToLower(directory)

	instance, loaded := sf_singletons.Load(directory)
	if loaded {
		return instance.(*VolatileFileSystem)
	}

	newInstance := NewVolatileFileSystem(mimeTypeDetection)
	actual, _ := sf_singletons.LoadOrStore(directory, newInstance)
	return actual.(*VolatileFileSystem)
}

// CreateDirectory implements IFileSystem.
func (fs *VolatileFileSystem) CreateDirectory(ctx context.Context, volume string, relPath string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	err = fs.CreateVolume(ctx, volume)
	if err != nil {
		return err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return err
	}
	path := JoinPaths(relPath, "")
	if _, ok := fs._volumes[volume]; !ok {
		fs._volumes[volume] = make(map[string][]byte)
	}
	fs._volumes[volume][path] = make([]byte, 0)
	return nil
}

// CreateVolume implements IFileSystem.
func (fs *VolatileFileSystem) CreateVolume(ctx context.Context, volume string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	if _, ok := fs._volumes[volume]; !ok {
		fs._volumes[volume] = make(map[string][]byte)
	}

	return nil
}

// DeleteDirectory implements IFileSystem.
func (fs *VolatileFileSystem) DeleteDirectory(ctx context.Context, volume string, relPath string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return err
	}

	if _, ok := fs._volumes[volume]; ok {
		files, err := fs.GetAllFileNames(ctx, volume, relPath)
		if err != nil {
			return err
		}
		for _, fileName := range files {
			path := JoinPaths(relPath, fileName)
			delete(fs._volumes[volume], path)
		}
		var dirPath = JoinPaths(relPath, "")
		delete(fs._volumes[volume], dirPath)
	}

	return nil
}

// DeleteFile implements IFileSystem.
func (fs *VolatileFileSystem) DeleteFile(ctx context.Context, volume string, relPath string, fileName string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return err
	}
	var path = JoinPaths(relPath, fileName)
	if _, ok := fs._volumes[volume]; ok {
		delete(fs._volumes[volume], path)
	}

	return nil
}

// DeleteVolume implements IFileSystem.
func (fs *VolatileFileSystem) DeleteVolume(ctx context.Context, volume string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}

	delete(fs._volumes, volume)
	return nil
}

// FileExists implements IFileSystem.
func (fs *VolatileFileSystem) FileExists(ctx context.Context, volume string, relPath string, fileName string) bool {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	volume, err := validateVolumeName(volume)
	if err != nil {
		return false
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return false
	}
	var path = JoinPaths(relPath, fileName)
	if strings.HasSuffix(path, "/") {
		return false
	}
	if _, ok := fs._volumes[volume]; ok {
		if _, ok := fs._volumes[volume][path]; ok {
			return true
		}
	}

	return false
}

// GetAllFileNames implements IFileSystem.
func (fs *VolatileFileSystem) GetAllFileNames(ctx context.Context, volume string, relPath string) ([]string, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	volume, err := validateVolumeName(volume)
	if err != nil {
		return nil, err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return nil, err
	}

	result := []string{}
	if _, ok := fs._volumes[volume]; ok {
		var path = JoinPaths(relPath, "")
		for key := range fs._volumes[volume] {
			if strings.HasPrefix(key, path) && key != path && !strings.Contains(key[len(path):], "/") {
				result = append(result, key[len(path):])
			}
		}
	}

	return result, nil
}

// ListVolumes implements IFileSystem.
func (fs *VolatileFileSystem) ListVolumes(ctx context.Context) ([]string, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	keys2 := make([]string, 0, len(fs._volumes))
	for k := range fs._volumes {
		keys2 = append(keys2, k)
	}

	return keys2, nil
}

// ReadAllFilesAsText implements IFileSystem.
func (fs *VolatileFileSystem) ReadAllFilesAsText(ctx context.Context, volume string, relPath string) (map[string]string, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	volume, err := validateVolumeName(volume)
	if err != nil {
		return nil, err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	if _, ok := fs._volumes[volume]; ok {
		var path = JoinPaths(relPath, "")
		for k, v := range fs._volumes[volume] {
			if strings.HasPrefix(k, path) {
				result[k] = string(v)
			}
		}
	}
	return result, nil
}

// ReadFileAsBinary implements IFileSystem.
func (fs *VolatileFileSystem) ReadFileAsBinary(ctx context.Context, volume string, relPath string, fileName string) ([]byte, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	volume, err := validateVolumeName(volume)
	if err != nil {
		return nil, err
	}
	if _, ok := fs._volumes[volume]; ok {

		fileName, err = validateFileName(fileName)
		if err != nil {
			return nil, err
		}
		var filePath = JoinPaths(relPath, fileName)

		if data, ok := fs._volumes[volume][filePath]; ok {
			return data, nil
		}
	}

	return nil, fmt.Errorf("file not found:%s %s", volume, relPath)
}

// ReadFileAsText implements IFileSystem.
func (fs *VolatileFileSystem) ReadFileAsText(ctx context.Context, volume string, relPath string, fileName string) (*string, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	data, err := fs.ReadFileAsBinary(ctx, volume, relPath, fileName)
	if err != nil {
		return nil, err
	}
	d := string(data)
	return &d, nil
}

// ReadFileInfo implements IFileSystem.
func (fs *VolatileFileSystem) ReadFileInfo(ctx context.Context, volume string, relPath string, fileName string) (*models.StreamableFileContent, error) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	volume, err := validateVolumeName(volume)
	if err != nil {
		return nil, err
	}
	if _, ok := fs._volumes[volume]; ok {
		relPath, err = validatePath(relPath)
		if err != nil {
			return nil, err
		}
		fileName, err = validateFileName(fileName)
		if err != nil {
			return nil, err
		}
		var filePath = JoinPaths(relPath, fileName)

		if data, ok := fs._volumes[volume][filePath]; ok {
			fileTypeP, err := fs.mimeTypeDetection.GetFileType(fileName)
			fileType := "application/octet-stream"
			if err == nil {
				fileType = *fileTypeP
			}
			result := models.NewStreamableFileContent(
				fileName,
				(int64)(len(data)),
				fileType,
				time.Now(),
				func(ctx context.Context) (io.ReadCloser, error) {
					return io.NopCloser(bytes.NewReader(data)), nil
				},
			)
			return result, nil
		}
	}

	return nil, fmt.Errorf("file not found:%s %s", volume, relPath)
}

// VolumeExists implements IFileSystem.
func (fs *VolatileFileSystem) VolumeExists(ctx context.Context, volume string) bool {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	volume, err := validateVolumeName(volume)
	if err != nil {
		return false
	}
	if _, ok := fs._volumes[volume]; ok {
		return true
	}

	return false
}

// WriteFile implements IFileSystem.
func (fs *VolatileFileSystem) WriteFile(ctx context.Context, volume string, relPath string, fileName string, streamContent io.Reader) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	if _, ok := fs._volumes[volume]; ok {
		relPath, err = validatePath(relPath)
		if err != nil {
			return err
		}
		fileName, err = validateFileName(fileName)
		if err != nil {
			return err
		}
		var path = JoinPaths(relPath, fileName)
		data, err := io.ReadAll(streamContent)
		if err != nil {
			return err
		}
		fs._volumes[volume][path] = data
	}

	return nil
}

// WriteFileAsText implements IFileSystem.
func (fs *VolatileFileSystem) WriteFileAsText(ctx context.Context, volume string, relPath string, fileName string, data string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	if _, ok := fs._volumes[volume]; ok {
		relPath, err = validatePath(relPath)
		if err != nil {
			return err
		}
		fileName, err = validateFileName(fileName)
		if err != nil {
			return err
		}
		var path = JoinPaths(relPath, fileName)
		fs._volumes[volume][path] = []byte(data)
	}

	return nil
}

func JoinPaths(a, b string) string {
	return strings.Trim(a, "/\\") + "/" + strings.Trim(b, "/\\")
}
