package domain

// IEventSourcing extends IAggregate with event sourcing specific methods, including versioning.
type IEventSourcing interface {
	IAggregateRoot
	IDomainEventHolder[IDomainEvent]
	Apply(event IDomainEvent) error
	AggregateVersion() int
}

type AggregateWithEventSourcing struct {
	Aggregate
	domainEvents []IDomainEvent
	Version      int
}

func (b *AggregateWithEventSourcing) AddDomainEvent(event IDomainEvent) {
	b.domainEvents = append(b.domainEvents, event)
}

func (b *AggregateWithEventSourcing) ClearDomainEvents() {
	b.domainEvents = nil
}

func (b *AggregateWithEventSourcing) DomainEvents() []IDomainEvent {
	return b.domainEvents
}

func (b *AggregateWithEventSourcing) AggregateVersion() int {
	return b.Version
}
