package infrastructure

import (
	"context"
	"fmt"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/options"
)

type ISnapshotStore[EventSourcing domain.IEventSourcing] interface {
	LoadSnapshot(ctx context.Context, id string) ([]EventSourcing, error)
	SaveSnapshot(ctx context.Context, aggregate EventSourcing) error
}

type ISnapshotStoreAsync[EventSourcing domain.IEventSourcing] interface {
	//  // LoadSnapshotAsync mothed usage
	//  func Stream(ctx context.Context, id string) error {
	//  	resultChan, errorChan := LoadSnapshotAsync(ctx, id)
	//  	select {
	//  	case datas := <-resultChan:
	//  	// handle data
	//  	case err := <-errorChan:
	//  	// handle error
	//  	case <-ctx.Done():
	//  	// handle timeout
	//  	}
	//  }
	//
	LoadSnapshotAsync(ctx context.Context, id string) (<-chan []EventSourcing, <-chan error)
	//  // SaveSnapshotAsync mothed usage
	//  func Stream(ctx context.Context, aggregate EventSourcing) error {
	//  	errorChan := SaveSnapshotAsync(ctx, aggregate)
	//  	select {
	//  	case err := <-errorChan:
	//  	// handle error
	//  	case <-ctx.Done():
	//  	// handle timeout
	//  	}
	//  }
	//
	SaveSnapshotAsync(ctx context.Context, aggregate EventSourcing) <-chan error
}

type SnapshotStoreRegistry struct {
	Options     *options.Options
	events      map[string]func(options.Options) ISnapshotStore[domain.IEventSourcing]
	eventAsyncs map[string]func(options.Options) ISnapshotStoreAsync[domain.IEventSourcing]
}

var DefaultSnapshotStoreRegistry *SnapshotStoreRegistry = NewSnapshotStoreRegistry()

func NewSnapshotStoreRegistry() *SnapshotStoreRegistry {
	return &SnapshotStoreRegistry{
		Options:     &options.Options{},
		events:      map[string]func(options.Options) ISnapshotStore[domain.IEventSourcing]{},
		eventAsyncs: map[string]func(options.Options) ISnapshotStoreAsync[domain.IEventSourcing]{},
	}
}

func (s *SnapshotStoreRegistry) RegisterComponent(componentFactory func(options.Options) ISnapshotStore[domain.IEventSourcing], componentAsyncFactory func(options.Options) ISnapshotStoreAsync[domain.IEventSourcing], names ...string) {
	for _, name := range names {
		s.events[fmt.Sprintf("snapshot-store-%s", name)] = componentFactory
		s.eventAsyncs[fmt.Sprintf("snapshot-store-async-%s", name)] = componentAsyncFactory
	}
}

func (s *SnapshotStoreRegistry) Create() (ISnapshotStore[domain.IEventSourcing], ISnapshotStoreAsync[domain.IEventSourcing], error) {
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
