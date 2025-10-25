package domain

import (
	"context"
)

type Repository[Aggregate AggregateRoot] interface {
	Get(ctx context.Context, id string) (*Aggregate, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, aggregate Aggregate) error
	Insert(ctx context.Context, aggregate Aggregate) error
}
