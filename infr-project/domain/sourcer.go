package domain

type IEventApplier[E IDomainEvent, R IEventSourcing] interface {
	Apply(aggregate R, event E) (R, error)
}

type IVersionManager[R IEventSourcing] interface {
	GetAllVersions(id string) ([]R, error)
	GetSpecificVersion(id string, version int) (*R, error)
	GetLatestVersion(id string) (*R, error)
}

type IEventSourcer[E IDomainEvent, R IEventSourcing] interface {
	IEventStore[E]
	IEventApplier[E, R]
	IVersionManager[R]
}
