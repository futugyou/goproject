package core

import (
	"context"
)

type IRepository[E IEntity, K any] interface {
	Insert(ctx context.Context, obj E) error
	Update(ctx context.Context, obj E, id K) error
	Delete(ctx context.Context, obj E, id K) error
	GetAll(ctx context.Context) ([]*E, error)
	Get(ctx context.Context, id K) (*E, error)
}
