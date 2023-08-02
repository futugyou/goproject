package core

import (
	"context"
)

type IRepository[E IEntity, K any] interface {
	Insert(ctx context.Context, obj E) error
	Update(ctx context.Context, obj E, id K) error
	Delete(ctx context.Context, id K) error
	GetAll(ctx context.Context) ([]*E, error)
	Get(ctx context.Context, id K) (*E, error)
	InsertMany(ctx context.Context, items []E) error
	Paging(ctx context.Context, page Paging) ([]*E, error)
}
