package domain

type IEventApplier[Event IDomainEvent, EventSourcing IEventSourcing] interface {
	Apply(aggregate EventSourcing, event Event) error
}

type IAggregateRetriever[Aggregate IAggregate] interface {
	RetrieveAllVersions(id string) ([]Aggregate, error)
	RetrieveSpecificVersion(id string, version int) (*Aggregate, error)
	RetrieveLatestVersion(id string) (*Aggregate, error)
}
