package infrastructure

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type ISnapshotStore[EventSourcing domain.IEventSourcing] interface {
	LoadSnapshot(id string) ([]EventSourcing, error)
	// LoadLatestSnapshot(id string) (*EventSourcing, error)
	SaveSnapshot(ctx context.Context, aggregate EventSourcing) error
}
