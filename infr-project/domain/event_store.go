package domain

import (
	"context"
)

type IEventStore[Event IDomainEvent] interface {
	Save(ctx context.Context, events []Event) error
	Load(ctx context.Context, id string) ([]Event, error)
	LoadGreaterthanVersion(ctx context.Context, id string, version int) ([]Event, error)
}

type IEventStoreAsync[Event IDomainEvent] interface {
	SaveAsync(ctx context.Context, events []Event) <-chan error
	LoadAsync(ctx context.Context, id string) (<-chan []Event, <-chan error)
	LoadGreaterthanVersionAsync(ctx context.Context, id string, version int) (<-chan []Event, <-chan error)
}
