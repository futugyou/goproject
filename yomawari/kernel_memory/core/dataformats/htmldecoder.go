package dataformats

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/dataformats"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/text"
)

type HtmlDecoder struct {
}

// Decode implements dataformats.IContentDecoder.
func (m *HtmlDecoder) Decode(ctx context.Context, fileName string) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("HtmlDecoder is nil")
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return m.DecodeBytes(ctx, content)
}

// DecodeBytes implements dataformats.IContentDecoder.
func (m *HtmlDecoder) DecodeBytes(ctx context.Context, content []byte) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("HtmlDecoder is nil")
	}
	return m.DecodeStream(ctx, bytes.NewReader(content))
}

// DecodeStream implements dataformats.IContentDecoder.
func (h *HtmlDecoder) DecodeStream(ctx context.Context, stream io.Reader) (*dataformats.FileContent, error) {
	if h == nil {
		return nil, fmt.Errorf("HtmlDecoder is nil")
	}
	doc, err := goquery.NewDocumentFromReader(stream)
	if err != nil {
		return nil, err
	}

	var result = &dataformats.FileContent{Sections: make([]dataformats.Chunk, 0), MimeType: pipeline.MimeTypes_PlainText}
	sentencesAreComplete := true
	chunk := dataformats.Chunk{Content: text.NormalizeNewlines(doc.Text(), true), Number: 1, Metadata: dataformats.ChunkMeta(&sentencesAreComplete, nil)}
	result.Sections = append(result.Sections, chunk)
	return result, nil
}

// SupportsMimeType implements dataformats.IContentDecoder.
func (h *HtmlDecoder) SupportsMimeType(ctx context.Context, mimeType string) bool {
	return strings.HasPrefix(mimeType, pipeline.MimeTypes_PlainText)
}

var _ dataformats.IContentDecoder = (*HtmlDecoder)(nil)
