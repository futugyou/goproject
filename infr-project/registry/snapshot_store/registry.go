package snapshot_store

import (
	"context"
	"fmt"

	"github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/infrastructure"
	"github.com/futugyou/infr-project/registry/options"
)

type Registry struct {
	Options *options.Options
	events  map[string]func(context.Context, options.Options) infrastructure.SnapshotStore[domain.EventSourcing]
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		Options: &options.Options{},
		events:  map[string]func(context.Context, options.Options) infrastructure.SnapshotStore[domain.EventSourcing]{},
	}
}

func (s *Registry) RegisterComponent(
	componentFactory func(context.Context, options.Options) infrastructure.SnapshotStore[domain.EventSourcing],
	names ...string,
) {
	for _, name := range names {
		s.events[fmt.Sprintf("snapshot-store-%s", name)] = componentFactory
	}
}

func (s *Registry) Create(ctx context.Context) (infrastructure.SnapshotStore[domain.EventSourcing], error) {
	if s.Options == nil {
		return nil, fmt.Errorf("options is nil")
	}

	name := s.Options.StoreType
	if method, ok := s.events[fmt.Sprintf("snapshot-store-%s", name)]; ok {
		return method(ctx, *s.Options), nil
	}

	return nil, fmt.Errorf("couldn't find event store %s", name)
}
