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

type MarkDownDecoder struct {
}

// Decode implements dataformats.IContentDecoder.
func (m *MarkDownDecoder) Decode(ctx context.Context, fileName string) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MarkDownDecoder is nil")
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return m.DecodeBytes(ctx, content)
}

// DecodeBytes implements dataformats.IContentDecoder.
func (m *MarkDownDecoder) DecodeBytes(ctx context.Context, content []byte) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MarkDownDecoder is nil")
	}
	var result = &dataformats.FileContent{Sections: make([]dataformats.Chunk, 0), MimeType: pipeline.MimeTypes_MarkDown}
	sentencesAreComplete := true
	result.Sections = append(result.Sections, dataformats.Chunk{Content: string(content), Number: 1, Metadata: dataformats.ChunkMeta(&sentencesAreComplete, nil)})
	return result, nil
}

// DecodeStream implements dataformats.IContentDecoder.
func (m *MarkDownDecoder) DecodeStream(ctx context.Context, stream io.Reader) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MarkDownDecoder is nil")
	}

	data, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	var result = &dataformats.FileContent{Sections: make([]dataformats.Chunk, 0), MimeType: pipeline.MimeTypes_MarkDown}
	sentencesAreComplete := true
	result.Sections = append(result.Sections, dataformats.Chunk{Content: string(data), Number: 1, Metadata: dataformats.ChunkMeta(&sentencesAreComplete, nil)})
	return result, nil
}

// SupportsMimeType implements dataformats.IContentDecoder.
func (m *MarkDownDecoder) SupportsMimeType(ctx context.Context, mimeType string) bool {
	return strings.HasPrefix(mimeType, pipeline.MimeTypes_MarkDown)
}

var _ dataformats.IContentDecoder = (*MarkDownDecoder)(nil)
