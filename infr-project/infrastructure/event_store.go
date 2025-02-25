package infrastructure

import (
	"context"
	"fmt"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/options"
)

type IEventStore[Event domain.IDomainEvent] interface {
	Save(ctx context.Context, events []Event) error
	Load(ctx context.Context, id string) ([]Event, error)
	LoadGreaterthanVersion(ctx context.Context, id string, version int) ([]Event, error)
}

type IEventStoreAsync[Event domain.IDomainEvent] interface {
	SaveAsync(ctx context.Context, events []Event) <-chan error
	LoadAsync(ctx context.Context, id string) (<-chan []Event, <-chan error)
	LoadGreaterthanVersionAsync(ctx context.Context, id string, version int) (<-chan []Event, <-chan error)
}

type EventStoreRegistry struct {
	Options     *options.Options
	stores      map[string]func(options.Options) IEventStore[domain.IDomainEvent]
	storeAsyncs map[string]func(options.Options) IEventStoreAsync[domain.IDomainEvent]
}

var DefaultEventStoreRegistry *EventStoreRegistry = NewEventStoreRegistry()

func NewEventStoreRegistry() *EventStoreRegistry {
	return &EventStoreRegistry{
		Options:     &options.Options{},
		stores:      map[string]func(options.Options) IEventStore[domain.IDomainEvent]{},
		storeAsyncs: map[string]func(options.Options) IEventStoreAsync[domain.IDomainEvent]{},
	}
}

func (s *EventStoreRegistry) RegisterComponent(componentFactory func(options.Options) IEventStore[domain.IDomainEvent], componentAsyncFactory func(options.Options) IEventStoreAsync[domain.IDomainEvent], names ...string) {
	for _, name := range names {
		s.stores[fmt.Sprintf("event-store-%s", name)] = componentFactory
		s.storeAsyncs[fmt.Sprintf("event-store-async-%s", name)] = componentAsyncFactory
	}
}

func (s *EventStoreRegistry) Create() (IEventStore[domain.IDomainEvent], IEventStoreAsync[domain.IDomainEvent], error) {
	if s.Options == nil {
		return nil, nil, fmt.Errorf("options is nil")
	}

	name := s.Options.StoreType
	if method, ok := s.stores[fmt.Sprintf("event-store-%s", name)]; ok {
		if methodasync, ok2 := s.storeAsyncs[fmt.Sprintf("event-store-async-%s", name)]; ok2 {
			return method(*s.Options), methodasync(*s.Options), nil
		}
	}

	return nil, nil, fmt.Errorf("couldn't find event store %s", name)
}
