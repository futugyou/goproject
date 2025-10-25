package infrastructure

import (
	"context"

	"github.com/futugyou/domaincore/domain"
)

type SnapshotStore[ES domain.EventSourcing] interface {
	LoadSnapshot(ctx context.Context, id string) ([]ES, error)
	SaveSnapshot(ctx context.Context, aggregate ES) error
}
