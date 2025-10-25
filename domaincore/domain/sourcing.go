package domain

// IEventSourcing extends IAggregate with event sourcing specific methods, including versioning.
type EventSourcing interface {
	DomainEventHolder[DomainEvent]
	Apply(event DomainEvent) error
	Replay(events []DomainEvent) error
	AggregateVersion() int
	Clone() EventSourcing
}

type AggregateWithEventSourcing struct {
	Aggregate
	domainEvents []DomainEvent
	Version      int
}

func (b *AggregateWithEventSourcing) AddDomainEvent(event DomainEvent) {
	b.domainEvents = append(b.domainEvents, event)
}

func (b *AggregateWithEventSourcing) ClearDomainEvents() {
	b.domainEvents = nil
}

func (b *AggregateWithEventSourcing) DomainEvents() []DomainEvent {
	// Return a copy to avoid external modification of the original slice.
	// If the following method has too much impact on performance, use return b.domainEvents directly.
	eventsCopy := make([]DomainEvent, len(b.domainEvents))
	copy(eventsCopy, b.domainEvents)
	return eventsCopy
}

func (b *AggregateWithEventSourcing) AggregateVersion() int {
	return b.Version
}
