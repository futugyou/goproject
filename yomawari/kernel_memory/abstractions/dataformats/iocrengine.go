package dataformats

import (
	"context"
	"io"
)

type IOcrEngine interface {
	ExtractTextFromImage(ctx context.Context, stream io.Reader) (*string, error)
}
