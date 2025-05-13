package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// FileCollection organizes a set of files and streams, ensures unique file names, and prepares them for upload
type FileCollection struct {
	mu        sync.Mutex
	filePaths map[string]string        // Original file path -> unique file name
	streams   map[string]io.ReadCloser // Unique file name -> associated stream
	fileNames map[string]struct{}      // Record the used file name to ensure uniqueness
}

// NewFileCollection creates a new FileCollection instance
func NewFileCollection() *FileCollection {
	return &FileCollection{
		filePaths: make(map[string]string),
		streams:   make(map[string]io.ReadCloser),
		fileNames: make(map[string]struct{}),
	}
}

// AddFile adds a local file
func (fc *FileCollection) AddFile(filePath string) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if _, exists := fc.filePaths[filePath]; exists {
		return nil // The file already exists, skip
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: '%s'", filePath)
	}

	fileName := getUniqueFileName(fc.fileNames, filePath)
	fc.filePaths[filePath] = fileName
	fc.fileNames[fileName] = struct{}{}
	return nil
}

// AddStream adds a stream
func (fc *FileCollection) AddStream(fileName string, content io.ReadCloser) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if content == nil {
		return
	}

	if strings.TrimSpace(fileName) == "" {
		fileName = "content.txt"
	}

	fileName = getUniqueFileName(fc.fileNames, fileName)
	fc.streams[fileName] = content
	fc.fileNames[fileName] = struct{}{}
}

// GetStreams gets all files and streams
func (fc *FileCollection) GetStreams() []struct {
	Name    string
	Content io.ReadCloser
} {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	var result []struct {
		Name    string
		Content io.ReadCloser
	}

	for filePath, fileName := range fc.filePaths {
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Ignore errors
		}
		result = append(result, struct {
			Name    string
			Content io.ReadCloser
		}{fileName, io.NopCloser(bytes.NewReader(data))})
	}

	for fileName, stream := range fc.streams {
		result = append(result, struct {
			Name    string
			Content io.ReadCloser
		}{fileName, stream})
	}

	return result
}

// Calculate unique file names
func getUniqueFileName(fileNames map[string]struct{}, filePath string) string {
	fileName := getBaseName(filePath)

	if _, exists := fileNames[fileName]; exists {
		dirHash := CalculateSHA256(getDirName(filePath))
		count := 0
		for {
			newName := fmt.Sprintf("%s%d_%s", dirHash, count, fileName)
			if _, exists := fileNames[newName]; !exists {
				return newName
			}
			count++
		}
	}

	return fileName
}

// Get the file name (without path)
func getBaseName(filePath string) string {
	parts := strings.Split(filePath, "/")
	return parts[len(parts)-1]
}

// Get the directory name
func getDirName(filePath string) string {
	index := strings.LastIndex(filePath, "/")
	if index == -1 {
		return ""
	}
	return filePath[:index]
}

// Calculate the SHA256 hash
func CalculateSHA256(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}
