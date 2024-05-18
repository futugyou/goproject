package infrastructure

import (
	"github.com/futugyou/infr-project/domain"
)

type ISnapshotStore[EventSourcing domain.IEventSourcing] interface {
	LoadSnapshot(id string) ([]EventSourcing, error)
	SaveSnapshot(aggregate EventSourcing) error
}
