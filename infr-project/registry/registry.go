package registry

import (
	"github.com/futugyou/infr-project/registry/event_store"
	"github.com/futugyou/infr-project/registry/publisher" 
	"github.com/futugyou/infr-project/registry/snapshot_store"
	"github.com/futugyou/infr-project/registry/unit"
)

type Registry struct {
	unitOfWork    *unit.Registry
	snapshotStore *snapshot_store.Registry 
	publisher     *publisher.Registry
	eventStore    *event_store.Registry
}

func New() *Registry {
	return &Registry{
		unitOfWork:    unit.DefaultRegistry,
		snapshotStore: snapshot_store.DefaultRegistry, 
		publisher:     publisher.DefaultRegistry,
		eventStore:    event_store.DefaultRegistry,
	}
}

func (o *Registry) WithUnitOfWork(registry *unit.Registry) *Registry {
	o.unitOfWork = registry
	return o
}

func (r *Registry) UnitOfWork() *unit.Registry {
	return r.unitOfWork
}

func (o *Registry) WithSnapshotStore(registry *snapshot_store.Registry) *Registry {
	o.snapshotStore = registry
	return o
}

func (r *Registry) SnapshotStore() *snapshot_store.Registry {
	return r.snapshotStore
}
 

func (o *Registry) WithPublisher(registry *publisher.Registry) *Registry {
	o.publisher = registry
	return o
}

func (r *Registry) Publisher() *publisher.Registry {
	return r.publisher
}

func (o *Registry) WithEventStore(registry *event_store.Registry) *Registry {
	o.eventStore = registry
	return o
}

func (r *Registry) EventStore() *event_store.Registry {
	return r.eventStore
}
