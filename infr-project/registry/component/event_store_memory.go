package component

import (
	"context"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/infrastructure_memory"
	"github.com/futugyou/infr-project/options"
	"github.com/futugyou/infr-project/registry/event_store"
)

func init() {
	event_store.DefaultRegistry.RegisterComponent(func(ctx context.Context, option options.Options) infrastructure.IEventStore[domain.IDomainEvent] {
		return infrastructure_memory.NewMemoryEventStore[domain.IDomainEvent]()
	}, func(ctx context.Context, option options.Options) infrastructure.IEventStore[domain.IDomainEvent] {
		return infrastructure_memory.NewMemoryEventStore[domain.IDomainEvent]()
	}, "memory")
}
