package dataformats

import (
	"context"
	"io"
)

type IContentDecoder interface {
	SupportsMimeType(ctx context.Context, mimeType string) bool
	Decode(ctx context.Context, fileName string) (*FileContent, error)
	DecodeBytes(ctx context.Context, bytes []byte) (*FileContent, error)
	DecodeStream(ctx context.Context, stream io.Reader) (*FileContent, error)
}
