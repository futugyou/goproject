package core

import "context"

type IDistributedCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	Refresh(ctx context.Context, key string) error
	Remove(ctx context.Context, key string) error
}
