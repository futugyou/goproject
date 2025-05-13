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
	"github.com/futugyou/yomawari/kernel_memory/abstractions/text"
	"github.com/nguyenthenguyen/docx"
)

// TODO: A more suitable library is needed
type MsWordDecoder struct {
}

// Decode implements dataformats.IContentDecoder.
func (m *MsWordDecoder) Decode(ctx context.Context, fileName string) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MsWordDecoder is nil")
	}
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return m.DecodeBytes(ctx, content)
}

// DecodeBytes implements dataformats.IContentDecoder.
func (m *MsWordDecoder) DecodeBytes(ctx context.Context, content []byte) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MsWordDecoder is nil")
	}
	return m.DecodeStream(ctx, bytes.NewReader(content))
}

// DecodeStream implements dataformats.IContentDecoder.
func (m *MsWordDecoder) DecodeStream(ctx context.Context, stream io.Reader) (*dataformats.FileContent, error) {
	if m == nil {
		return nil, fmt.Errorf("MsWordDecoder is nil")
	}

	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, stream)
	if err != nil {
		return nil, err
	}

	readerAt := bytes.NewReader(buff.Bytes())
	r, err := docx.ReadDocxFromMemory(readerAt, size)
	if err != nil {
		return nil, err
	}

	doc := r.Editable().GetContent()
	var result = &dataformats.FileContent{Sections: make([]dataformats.Chunk, 0), MimeType: pipeline.MimeTypes_MsWordX}
	sentencesAreComplete := true
	chunk := dataformats.Chunk{Content: text.NormalizeNewlines(doc, true), Number: 1, Metadata: dataformats.ChunkMeta(&sentencesAreComplete, nil)}
	result.Sections = append(result.Sections, chunk)
	return result, nil
}

// SupportsMimeType implements dataformats.IContentDecoder.
func (m *MsWordDecoder) SupportsMimeType(ctx context.Context, mimeType string) bool {
	return strings.HasPrefix(mimeType, pipeline.MimeTypes_MsWordX)
}

var _ dataformats.IContentDecoder = (*MsWordDecoder)(nil)
