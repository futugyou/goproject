package eventsourcing

type ISnapshotter[E IEvent, R IAggregate] interface {
	TakeSnapshot(aggregate R) error
	RestoreFromSnapshot(id string) (*R, error)
	RestoreFromSnapshotByVersion(id string, version int) (*R, error)
}
