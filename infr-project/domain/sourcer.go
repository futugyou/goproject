package domain

type IEventApplier[Event IDomainEvent, EventSourcing IEventSourcing] interface {
	Apply(aggregate EventSourcing, event Event) (EventSourcing, error)
}

type IVersionManager[EventSourcing IEventSourcing] interface {
	GetAllVersions(id string) ([]EventSourcing, error)
	GetSpecificVersion(id string, version int) (*EventSourcing, error)
	GetLatestVersion(id string) (*EventSourcing, error)
}

type IEventSourcer[Event IDomainEvent, EventSourcing IEventSourcing] interface {
	IEventStore[Event]
	IEventApplier[Event, EventSourcing]
	IVersionManager[EventSourcing]
}
