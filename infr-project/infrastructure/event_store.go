package infrastructure

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type IEventStore[Event domain.IDomainEvent] interface {
	Save(ctx context.Context, events []Event) error
	Load(ctx context.Context, id string) ([]Event, error)
	LoadGreaterthanVersion(ctx context.Context, id string, version int) ([]Event, error)
}
