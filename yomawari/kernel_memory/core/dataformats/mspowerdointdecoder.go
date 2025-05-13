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

type MsPowerPointDecoder struct {
}

// Decode implements dataformats.IContentDecoder.
func (m *MsPowerPointDecoder) Decode(ctx context.Context, fileName string) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MsPowerPointDecoder is nil")
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return m.DecodeBytes(ctx, content)
}

// DecodeBytes implements dataformats.IContentDecoder.
func (m *MsPowerPointDecoder) DecodeBytes(ctx context.Context, content []byte) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MsPowerPointDecoder is nil")
	}
	return m.DecodeStream(ctx, bytes.NewReader(content))
}

// DecodeStream implements dataformats.IContentDecoder.
func (m *MsPowerPointDecoder) DecodeStream(ctx context.Context, stream io.Reader) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MsPowerPointDecoder is nil")
	}
	panic("unimplemented")
}

// SupportsMimeType implements dataformats.IContentDecoder.
func (m *MsPowerPointDecoder) SupportsMimeType(ctx context.Context, mimeType string) bool {
	return strings.HasPrefix(mimeType, pipeline.MimeTypes_MsPowerPointX)
}

var _ dataformats.IContentDecoder = (*MsPowerPointDecoder)(nil)
