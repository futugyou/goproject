package eventsourcing

type IEvent interface {
	EventType() string
}

type IAggregate interface {
	AggregateName() string
}

type IEventSourcer[E IEvent, R IAggregate] interface {
	Add(event E) error
	Save(events []E) error
	Load(id string) ([]E, error)
	Apply(aggregate R, event E) R
	GetAlltVersions() []R
	GetSpecificVersion(version int) R
}
