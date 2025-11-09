package service

import "context"

type Screenshot interface {
	Create(ctx context.Context, url string) (*string, error)
}
