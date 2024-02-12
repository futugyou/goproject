package core

import (
	"context"
)

type IRepository[E IEntity, K any] interface {
	GetAll(ctx context.Context) ([]E, error)
	Paging(ctx context.Context, page Paging) ([]E, error)
	InsertMany(ctx context.Context, items []E) error
}
