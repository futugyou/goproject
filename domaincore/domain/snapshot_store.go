package domain

import (
	"context"
)

type SnapshotStore[ES EventSourcing] interface {
	LoadSnapshot(ctx context.Context, id string) ([]ES, error)
	SaveSnapshot(ctx context.Context, aggregate ES) error
}
