package registry

import (
	"github.com/futugyou/infr-project/registry/event_store"
	"github.com/futugyou/infr-project/registry/publisher"
	"github.com/futugyou/infr-project/registry/screenshot"
	"github.com/futugyou/infr-project/registry/snapshot_store"
	"github.com/futugyou/infr-project/registry/unit"
)

type Registry struct {
	UnitOfWork    *unit.Registry
	SnapshotStore *snapshot_store.Registry
	Screenshot    *screenshot.Registry
	Publisher     *publisher.Registry
	EventStore    *event_store.Registry
}

func New() *Registry {
	return &Registry{
		UnitOfWork:    unit.DefaultRegistry,
		SnapshotStore: snapshot_store.DefaultRegistry,
		Screenshot:    screenshot.DefaultRegistry,
		Publisher:     publisher.DefaultRegistry,
		EventStore:    event_store.DefaultRegistry,
	}
}
