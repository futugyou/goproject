package domain

import (
	"context"
)

type IRepository[Aggregate IAggregateRoot] interface {
	Get(ctx context.Context, id string) (*Aggregate, error)
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
	Update(ctx context.Context, aggregate Aggregate) error
	Insert(ctx context.Context, aggregate Aggregate) error
}

type IRepositoryAsync[Aggregate IAggregateRoot] interface {
	GetAsync(ctx context.Context, id string) (<-chan *Aggregate, <-chan error)
	DeleteAsync(ctx context.Context, id string) <-chan error
	SoftDeleteAsync(ctx context.Context, id string) <-chan error
	UpdateAsync(ctx context.Context, aggregate Aggregate) <-chan error
	InsertAsync(ctx context.Context, aggregate Aggregate) <-chan error
}
