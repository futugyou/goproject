package event_store

import (
	"context"
	"fmt"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/options"
)

type Registry struct {
	Options *options.Options
	stores  map[string]func(context.Context, options.Options) infrastructure.IEventStore[domain.IDomainEvent]
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		Options: &options.Options{},
		stores:  map[string]func(context.Context, options.Options) infrastructure.IEventStore[domain.IDomainEvent]{},
	}
}

func (s *Registry) RegisterComponent(
	componentFactory func(context.Context, options.Options) infrastructure.IEventStore[domain.IDomainEvent],
	names ...string,
) {
	for _, name := range names {
		s.stores[fmt.Sprintf("event-store-%s", name)] = componentFactory
	}
}

func (s *Registry) Create(ctx context.Context) (infrastructure.IEventStore[domain.IDomainEvent], error) {
	if s.Options == nil {
		return nil, fmt.Errorf("options is nil")
	}

	name := s.Options.StoreType
	if method, ok := s.stores[fmt.Sprintf("event-store-%s", name)]; ok {
		return method(ctx, *s.Options), nil
	}

	return nil, fmt.Errorf("couldn't find event store %s", name)
}
