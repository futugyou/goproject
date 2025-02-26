package snapshot_store

import (
	"fmt"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/options"
)

type Registry struct {
	Options     *options.Options
	events      map[string]func(options.Options) infrastructure.ISnapshotStore[domain.IEventSourcing]
	eventAsyncs map[string]func(options.Options) infrastructure.ISnapshotStoreAsync[domain.IEventSourcing]
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		Options:     &options.Options{},
		events:      map[string]func(options.Options) infrastructure.ISnapshotStore[domain.IEventSourcing]{},
		eventAsyncs: map[string]func(options.Options) infrastructure.ISnapshotStoreAsync[domain.IEventSourcing]{},
	}
}

func (s *Registry) RegisterComponent(
	componentFactory func(options.Options) infrastructure.ISnapshotStore[domain.IEventSourcing],
	componentAsyncFactory func(options.Options) infrastructure.ISnapshotStoreAsync[domain.IEventSourcing],
	names ...string,
) {
	for _, name := range names {
		s.events[fmt.Sprintf("snapshot-store-%s", name)] = componentFactory
		s.eventAsyncs[fmt.Sprintf("snapshot-store-async-%s", name)] = componentAsyncFactory
	}
}

func (s *Registry) Create() (infrastructure.ISnapshotStore[domain.IEventSourcing], infrastructure.ISnapshotStoreAsync[domain.IEventSourcing], error) {
	if s.Options == nil {
		return nil, nil, fmt.Errorf("options is nil")
	}

	name := s.Options.StoreType
	if method, ok := s.events[fmt.Sprintf("snapshot-store-%s", name)]; ok {
		if methodasync, ok2 := s.eventAsyncs[fmt.Sprintf("snapshot-store-async-%s", name)]; ok2 {
			return method(*s.Options), methodasync(*s.Options), nil
		}
	}

	return nil, nil, fmt.Errorf("couldn't find event store %s", name)
}
