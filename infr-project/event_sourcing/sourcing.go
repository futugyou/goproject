package eventsourcing

type IEvent interface {
	EventType() string
}

type IAggregate interface {
	AggregateName() string
	Apply(event IEvent) error
}

type IEventSourcer[E IEvent, R IAggregate] interface {
	Add(event E) error
	Save(events []E) error
	Load(id string) ([]E, error)
	Apply(aggregate R, event E) R
	GetAllVersions(id string) ([]R, error)
	GetSpecificVersion(id string, version int) (*R, error)
}
