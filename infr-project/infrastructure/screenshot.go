package infrastructure

import "context"

type IScreenshot interface {
	Create(ctx context.Context, url string) (*string, error)
}
