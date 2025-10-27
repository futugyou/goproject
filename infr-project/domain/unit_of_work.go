package domain

import (
	"context"
)

type IUnitOfWork interface {
	BeginTransaction(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
