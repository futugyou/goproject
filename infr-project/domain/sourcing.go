package domain

// IEventSourcing extends IAggregate with event sourcing specific methods, including versioning.
type IEventSourcing interface {
	IAggregate
	// IDomainEventHolder
	Apply(event IDomainEvent) (IEventSourcing, error)
	AggregateVersion() int
}
