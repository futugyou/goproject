package domain

import (
	"context"
)

const DATA_NOT_FOUND_MESSAGE string = "data not found"

type Repository[Aggregate AggregateRoot] interface {
	FindByID(ctx context.Context, id string) (*Aggregate, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, aggregate Aggregate) error
	Insert(ctx context.Context, aggregate Aggregate) error
	Find(ctx context.Context, options *QueryOptions) ([]Aggregate, error)
}
