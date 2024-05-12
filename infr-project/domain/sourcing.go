package domain

// IEventSourcing extends IAggregate with event sourcing specific methods, including versioning.
type IEventSourcing interface {
	IAggregate
	IDomainEventHolder[IDomainEvent]
	Apply(event IDomainEvent) error
	AggregateVersion() int
}
