package infrastructure

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type ISnapshotStore[EventSourcing domain.IEventSourcing] interface {
	LoadSnapshot(ctx context.Context, id string) ([]EventSourcing, error)
	// LoadLatestSnapshot(ctx context.Context, id string) (*EventSourcing, error)
	SaveSnapshot(ctx context.Context, aggregate EventSourcing) error
}
