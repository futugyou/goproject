package domain

import (
	"context"
)

type ISnapshotStore[EventSourcing IEventSourcing] interface {
	LoadSnapshot(ctx context.Context, id string) ([]EventSourcing, error)
	SaveSnapshot(ctx context.Context, aggregate EventSourcing) error
}

type ISnapshotStoreAsync[EventSourcing IEventSourcing] interface {
	LoadSnapshotAsync(ctx context.Context, id string) (<-chan []EventSourcing, <-chan error)
	SaveSnapshotAsync(ctx context.Context, aggregate EventSourcing) <-chan error
}
