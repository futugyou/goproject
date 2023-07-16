package core

import (
	"context"
)

type IInsertRepository[E IEntity] interface {
	Insert(ctx context.Context, obj E) error
}

type IUpdateRepository[E IEntity, K any] interface {
	Update(ctx context.Context, obj E, id K) error
}

type IDeleteRepository[E IEntity, K any] interface {
	Delete(ctx context.Context, obj E, id K) error
}

type IGetAllRepository[E IEntity] interface {
	GetAll(ctx context.Context) ([]*E, error)
}

type IGetRepository[E IEntity, K any] interface {
	Get(ctx context.Context, id K) (*E, error)
}
