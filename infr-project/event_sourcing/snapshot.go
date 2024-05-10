package eventsourcing

type ISnapshotter interface {
	TakeSnapshot() error
	RestoreFromSnapshot() error
}
