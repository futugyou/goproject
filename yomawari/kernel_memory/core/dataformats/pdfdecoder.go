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

type PdfDecoder struct {
}

// Decode implements dataformats.IContentDecoder.
func (m *PdfDecoder) Decode(ctx context.Context, fileName string) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("PdfDecoder is nil")
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return m.DecodeBytes(ctx, content)
}

// DecodeBytes implements dataformats.IContentDecoder.
func (m *PdfDecoder) DecodeBytes(ctx context.Context, content []byte) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("PdfDecoder is nil")
	}
	return m.DecodeStream(ctx, bytes.NewReader(content))
}

// DecodeStream implements dataformats.IContentDecoder.
func (m *PdfDecoder) DecodeStream(ctx context.Context, stream io.Reader) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("PdfDecoder is nil")
	}
	panic("unimplemented")
}

// SupportsMimeType implements dataformats.IContentDecoder.
func (p *PdfDecoder) SupportsMimeType(ctx context.Context, mimeType string) bool {
	return strings.HasPrefix(mimeType, pipeline.MimeTypes_Pdf)
}

var _ dataformats.IContentDecoder = (*PdfDecoder)(nil)
