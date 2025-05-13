package dataformats

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/dataformats"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
)

type TextDecoder struct {
}

// Decode implements dataformats.IContentDecoder.
func (m *TextDecoder) Decode(ctx context.Context, fileName string) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("TextDecoder is nil")
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return m.DecodeBytes(ctx, content)
}

// DecodeBytes implements dataformats.IContentDecoder.
func (m *TextDecoder) DecodeBytes(ctx context.Context, content []byte) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("TextDecoder is nil")
	}
	var result = &dataformats.FileContent{Sections: make([]dataformats.Chunk, 0), MimeType: pipeline.MimeTypes_PlainText}
	sentencesAreComplete := true
	result.Sections = append(result.Sections, dataformats.Chunk{Content: string(content), Number: 1, Metadata: dataformats.ChunkMeta(&sentencesAreComplete, nil)})
	return result, nil
}

// DecodeStream implements dataformats.IContentDecoder.
func (m *TextDecoder) DecodeStream(ctx context.Context, stream io.Reader) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("TextDecoder is nil")
	}

	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	var result = &dataformats.FileContent{Sections: make([]dataformats.Chunk, 0), MimeType: pipeline.MimeTypes_PlainText}
	sentencesAreComplete := true
	result.Sections = append(result.Sections, dataformats.Chunk{Content: string(data), Number: 1, Metadata: dataformats.ChunkMeta(&sentencesAreComplete, nil)})
	return result, nil
}

// SupportsMimeType implements dataformats.IContentDecoder.
func (m *TextDecoder) SupportsMimeType(ctx context.Context, mimeType string) bool {
	return strings.HasPrefix(mimeType, pipeline.MimeTypes_PlainText) || strings.HasPrefix(mimeType, pipeline.MimeTypes_Json)
}

var _ dataformats.IContentDecoder = (*TextDecoder)(nil)
