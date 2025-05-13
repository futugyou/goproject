package filesystem

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
)

type DiskFileSystem struct {
	dataPath          string
	mutex             sync.Mutex
	mimeTypeDetection pipeline.IMimeTypeDetection
}

var invalidCharsRegex = regexp.MustCompile(`[\s|\\|/|\0|'|\"|:|;|,|~|!|?|*|+|=|^|@|#|$|%|&]`)

func NewDiskFileSystem(directory string, mimeTypeDetection pipeline.IMimeTypeDetection) *DiskFileSystem {
	if mimeTypeDetection == nil {
		mimeTypeDetection = &pipeline.MimeTypesDetection{}
	}
	fs := &DiskFileSystem{dataPath: directory, mimeTypeDetection: mimeTypeDetection}
	fs.createDirectory(directory)
	return fs
}

func (fs *DiskFileSystem) CreateVolume(ctx context.Context, volume string) error {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	return fs.createDirectory(filepath.Join(fs.dataPath, volume))
}

func (fs *DiskFileSystem) VolumeExists(ctx context.Context, volume string) bool {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(fs.dataPath, volume))
	return err == nil || !os.IsNotExist(err)
}

func (fs *DiskFileSystem) DeleteVolume(ctx context.Context, volume string) error {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(fs.dataPath, volume))
}

func (fs *DiskFileSystem) ListVolumes(ctx context.Context) ([]string, error) {
	dirs, err := os.ReadDir(fs.dataPath)
	if err != nil {
		return nil, err
	}
	var volumes []string
	for _, dir := range dirs {
		if dir.IsDir() {
			volumes = append(volumes, dir.Name())
		}
	}
	return volumes, nil
}

func (fs *DiskFileSystem) CreateDirectory(ctx context.Context, volume, relPath string) error {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return err
	}

	path := filepath.Join(fs.dataPath, volume, relPath)
	return fs.createDirectory(path)
}

func (fs *DiskFileSystem) DeleteDirectory(ctx context.Context, volume, relPath string) error {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return err
	}

	path := filepath.Join(fs.dataPath, volume, relPath)
	if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
		return os.Remove(path)
	}

	return nil
}

func (fs *DiskFileSystem) WriteFile(ctx context.Context, volume, relPath, fileName string, content io.Reader) error {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return err
	}

	fileName, err = validateFileName(fileName)
	if err != nil {
		return err
	}

	path := filepath.Join(fs.dataPath, volume, relPath)
	if err := fs.createDirectory(path); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	return err
}

func (fs *DiskFileSystem) WriteFileAsText(ctx context.Context, volume, relPath, fileName, data string) error {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return err
	}

	fileName, err = validateFileName(fileName)
	if err != nil {
		return err
	}

	path := filepath.Join(fs.dataPath, volume, relPath)
	if err := fs.createDirectory(path); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	return err
}

func (fs *DiskFileSystem) FileExists(ctx context.Context, volume, relPath, fileName string) bool {
	path := filepath.Join(fs.dataPath, volume, relPath, fileName)
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func (fs *DiskFileSystem) ReadFileAsBinary(ctx context.Context, volume, relPath, fileName string) ([]byte, error) {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return nil, err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return nil, err
	}

	fileName, err = validateFileName(fileName)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(fs.dataPath, volume, relPath, fileName)

	if !fs.FileExists(ctx, volume, relPath, fileName) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	return os.ReadFile(path)
}

func (fs *DiskFileSystem) ReadFileInfo(ctx context.Context, volume, relPath, fileName string) (*models.StreamableFileContent, error) {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return nil, err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return nil, err
	}

	fileName, err = validateFileName(fileName)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(fs.dataPath, volume, relPath, fileName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("directory not found: " + path)
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	fileTypeP, err := fs.mimeTypeDetection.GetFileType(fileName)
	fileType := "application/octet-stream"
	if err == nil {
		fileType = *fileTypeP
	}
	result := models.NewStreamableFileContent(
		fileName,
		fileSize,
		fileType,
		fileInfo.ModTime(),
		func(ctx context.Context) (io.ReadCloser, error) {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}

			return io.NopCloser(bytes.NewReader(data)), nil
		},
	)
	return result, nil
}

func (fs *DiskFileSystem) ReadFileAsText(ctx context.Context, volume, relPath, fileName string) (*string, error) {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return nil, err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return nil, err
	}

	fileName, err = validateFileName(fileName)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(fs.dataPath, volume, relPath, fileName)

	if !fs.FileExists(ctx, volume, relPath, fileName) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := string(data)
	return &result, nil
}

func (fs *DiskFileSystem) ReadAllFilesAsText(ctx context.Context, volume, relPath string) (map[string]string, error) {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return nil, err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(fs.dataPath, volume, relPath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("directory not found: " + path)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(path, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			result[file.Name()] = string(content)
		}
	}

	return result, nil
}

func (fs *DiskFileSystem) GetAllFileNames(ctx context.Context, volume, relPath string) ([]string, error) {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return nil, err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(fs.dataPath, volume, relPath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("directory not found: " + path)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, file := range files {
		if !file.IsDir() {
			result = append(result, file.Name())
		}
	}

	return result, nil
}

func (fs *DiskFileSystem) DeleteFile(ctx context.Context, volume, relPath, fileName string) error {
	volume, err := validateVolumeName(volume)
	if err != nil {
		return err
	}
	relPath, err = validatePath(relPath)
	if err != nil {
		return err
	}
	fileName, err = validateFileName(fileName)
	if err != nil {
		return err
	}
	path := filepath.Join(fs.dataPath, volume, relPath, fileName)
	if !fs.FileExists(ctx, volume, relPath, fileName) {
		return fmt.Errorf("file not found: %s", path)
	}
	return os.Remove(path)
}

func (fs *DiskFileSystem) createDirectory(path string) error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	return os.MkdirAll(path, 0755)
}

func validateVolumeName(volume string) (string, error) {
	if volume == "" {
		return "__default__", nil
	}
	if invalidCharsRegex.MatchString(volume) {
		return "", errors.New("invalid volume name")
	}
	return volume, nil
}

func validatePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	if matched, _ := regexp.MatchString(`[\\|:]`, path); matched {
		return "", errors.New("invalid path characters")
	}
	return path, nil
}

func validateFileName(fileName string) (string, error) {
	if fileName == "" {
		return "", errors.New("file name cannot be empty")
	}
	if matched, _ := regexp.MatchString(`[/\\|:]`, fileName); matched {
		return "", errors.New("invalid file name characters")
	}
	return fileName, nil
}
