package event_store

import (
	"context"
	"fmt"

	"github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/infrastructure"
	"github.com/futugyou/infr-project/registry/options"
)

type Registry struct {
	Options *options.Options
	stores  map[string]func(context.Context, options.Options) infrastructure.EventStore[domain.DomainEvent]
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		Options: &options.Options{},
		stores:  map[string]func(context.Context, options.Options) infrastructure.EventStore[domain.DomainEvent]{},
	}
}

func (s *Registry) RegisterComponent(
	componentFactory func(context.Context, options.Options) infrastructure.EventStore[domain.DomainEvent],
	names ...string,
) {
	for _, name := range names {
		s.stores[fmt.Sprintf("event-store-%s", name)] = componentFactory
	}
}

func (s *Registry) Create(ctx context.Context) (infrastructure.EventStore[domain.DomainEvent], error) {
	if s.Options == nil {
		return nil, fmt.Errorf("options is nil")
	}

	name := s.Options.StoreType
	if method, ok := s.stores[fmt.Sprintf("event-store-%s", name)]; ok {
		return method(ctx, *s.Options), nil
	}

	return nil, fmt.Errorf("couldn't find event store %s", name)
}
