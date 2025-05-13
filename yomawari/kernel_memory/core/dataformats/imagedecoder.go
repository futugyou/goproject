package dataformats

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/dataformats"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
)

type ImageDecoder struct {
	ocrEngine dataformats.IOcrEngine
}

// Decode implements dataformats.IContentDecoder.
func (m *ImageDecoder) Decode(ctx context.Context, fileName string) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("ImageDecoder is nil")
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return m.DecodeBytes(ctx, content)
}

// DecodeBytes implements dataformats.IContentDecoder.
func (m *ImageDecoder) DecodeBytes(ctx context.Context, content []byte) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("ImageDecoder is nil")
	}
	return m.DecodeStream(ctx, bytes.NewReader(content))
}

// DecodeStream implements dataformats.IContentDecoder.
func (i *ImageDecoder) DecodeStream(ctx context.Context, stream io.Reader) (*dataformats.FileContent, error) {
	if i == nil {
		return nil, fmt.Errorf("ImageDecoder is nil")
	}

	var result = &dataformats.FileContent{Sections: make([]dataformats.Chunk, 0), MimeType: pipeline.MimeTypes_PlainText}
	content, err := i.ImageToTextReader(ctx, stream)
	if err != nil {
		return nil, err
	}

	sentencesAreComplete := true
	result.Sections = append(result.Sections, dataformats.Chunk{Content: *content, Number: 1, Metadata: dataformats.ChunkMeta(&sentencesAreComplete, nil)})

	return result, nil
}

// SupportsMimeType implements dataformats.IContentDecoder.
func (i *ImageDecoder) SupportsMimeType(ctx context.Context, mimeType string) bool {
	return strings.HasPrefix(mimeType, pipeline.MimeTypes_ImageJpeg) ||
		strings.HasPrefix(mimeType, pipeline.MimeTypes_ImagePng) ||
		strings.HasPrefix(mimeType, pipeline.MimeTypes_ImageTiff)
}

func (i *ImageDecoder) ImageToText(ctx context.Context, fileName string) (*string, error) {
	if i.ocrEngine == nil {
		return nil, fmt.Errorf("ocrEngine is nil")
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return i.ImageToTextByte(ctx, content)
}

func (i *ImageDecoder) ImageToTextByte(ctx context.Context, content []byte) (*string, error) {
	if i.ocrEngine == nil {
		return nil, fmt.Errorf("ocrEngine is nil")
	}
	return i.ImageToTextReader(ctx, bytes.NewReader(content))
}

func (i *ImageDecoder) ImageToTextReader(ctx context.Context, content io.Reader) (*string, error) {
	if i.ocrEngine == nil {
		return nil, fmt.Errorf("ocrEngine is nil")
	}
	readCloser := io.NopCloser(content)
	return i.ImageToTextStream(ctx, readCloser)
}

func (i *ImageDecoder) ImageToTextStream(ctx context.Context, content io.ReadCloser) (*string, error) {
	if i.ocrEngine == nil {
		return nil, fmt.Errorf("ocrEngine is nil")
	}

	return i.ocrEngine.ExtractTextFromImage(ctx, content)
}

var _ dataformats.IContentDecoder = (*ImageDecoder)(nil)
