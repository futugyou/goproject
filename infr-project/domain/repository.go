package domain

import (
	"context"
)

// Deprecated: IEventSourcingRepository is deprecated, use ISnapshotStore and IEventStore.
type IEventSourcingRepository[EventSourcing IEventSourcing] interface {
	Load(ctx context.Context, id string) (*EventSourcing, error)
	LoadAll(ctx context.Context, id string) ([]EventSourcing, error)
	Save(ctx context.Context, aggregate EventSourcing) error
}

type IRepository[Aggregate IAggregateRoot] interface {
	Get(ctx context.Context, id string) (*Aggregate, error)
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
	Update(ctx context.Context, aggregate Aggregate) error
	Insert(ctx context.Context, aggregate Aggregate) error
}
