package domain

type ISnapshotter[Event IDomainEvent, EventSourcing IEventSourcing] interface {
	TakeSnapshot(aggregate EventSourcing) error
	RestoreFromSnapshot(id string) (*EventSourcing, error)
	RestoreFromSnapshotByVersion(id string, version int) (*EventSourcing, error)
}
