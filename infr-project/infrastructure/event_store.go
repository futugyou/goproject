package infrastructure

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type IEventStore[Event domain.IDomainEvent] interface {
	Save(ctx context.Context, events []Event) error
	Load(id string) ([]Event, error)
	LoadGreaterthanVersion(id string, version int) ([]Event, error)
}
