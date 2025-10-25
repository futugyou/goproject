package domain

type Snapshotter[E DomainEvent, ES EventSourcing] interface {
	TakeSnapshot(aggregate ES) error
	RestoreFromSnapshot(id string) (*ES, error)
	RestoreFromSnapshotByVersion(id string, version int) (*ES, error)
}
