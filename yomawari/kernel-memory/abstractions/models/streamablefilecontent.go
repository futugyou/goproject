package models

import (
	"bytes"
	"context"
	"io"
	"time"
)

type StreamableFileContent struct {
	stream      io.ReadCloser
	fileName    string
	fileSize    int64
	fileType    string
	lastWrite   time.Time
	initialized bool
	getStream   func(ctx context.Context) (io.ReadCloser, error)
}

func NewStreamableFileContent(fileName string, fileSize int64, fileType string, lastWrite time.Time, getStream func(ctx context.Context) (io.ReadCloser, error)) *StreamableFileContent {
	if len(fileType) == 0 {
		fileType = "application/octet-stream"
	}

	content := &StreamableFileContent{
		fileName:  fileName,
		fileSize:  fileSize,
		fileType:  fileType,
		lastWrite: lastWrite,
	}

	if getStream == nil {
		getStream = func(ctx context.Context) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		}
	}

	content.getStream = func(ctx context.Context) (io.ReadCloser, error) {
		if !content.initialized {
			var err error
			content.stream, err = getStream(ctx)
			if err != nil {
				return nil, err
			}
			content.initialized = true
		}
		return content.stream, nil
	}

	return content
}

func (content *StreamableFileContent) GetStream() func(ctx context.Context) (io.ReadCloser, error) {
	if content == nil || content.getStream == nil {
		return func(ctx context.Context) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		}
	}
	return content.getStream
}

func (content *StreamableFileContent) FileName() string {
	return content.fileName
}

func (content *StreamableFileContent) FileType() string {
	return content.fileType
}

func (content *StreamableFileContent) FileSize() int64 {
	return content.fileSize
}

func (content *StreamableFileContent) LastWrite() time.Time {
	return content.lastWrite
}

func (s *StreamableFileContent) Close() error {
	if s.stream != nil {
		err := s.stream.Close()
		s.stream = nil
		return err
	}
	return nil
}
