package component

import (
	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/infrastructure_memory"
	"github.com/futugyou/infr-project/options"
	"github.com/futugyou/infr-project/registry/snapshot_store"
)

func init() {
	snapshot_store.DefaultRegistry.RegisterComponent(func(option options.Options) infrastructure.ISnapshotStore[domain.IEventSourcing] {
		return infrastructure_memory.NewMemorySnapshotStore[domain.IEventSourcing]()
	}, func(option options.Options) infrastructure.ISnapshotStoreAsync[domain.IEventSourcing] {
		return infrastructure_memory.NewMemorySnapshotStore[domain.IEventSourcing]()
	}, "memory")
}
