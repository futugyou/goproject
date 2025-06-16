package domain

// IEventSourcing extends IAggregate with event sourcing specific methods, including versioning.
type IEventSourcing interface {
	IDomainEventHolder[IDomainEvent]
	Apply(event IDomainEvent) error
	Replay(events []IDomainEvent) error
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
	// Return a copy to avoid external modification of the original slice.
	// If the following method has too much impact on performance, use return b.domainEvents directly.
	eventsCopy := make([]IDomainEvent, len(b.domainEvents))
	copy(eventsCopy, b.domainEvents)
	return eventsCopy
}

func (b *AggregateWithEventSourcing) AggregateVersion() int {
	return b.Version
}
