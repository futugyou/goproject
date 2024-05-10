package eventsourcing

type IEvent interface {
	EventType() string
}

type IAggregate interface {
	AggregateName() string
	AggregateId() string
	Apply(event IEvent) error
}

type IEventSourcer[E IEvent, R IAggregate] interface {
	Save(events []E) error
	Load(id string) ([]E, error)
	Apply(aggregate R, event E) R
	GetAllVersions(aggregate R) ([]R, error)
	GetSpecificVersion(aggregate R, version int) (*R, error)
}
