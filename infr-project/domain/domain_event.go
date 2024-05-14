package domain

// IEvent represents the interface for events.
type IDomainEvent interface {
	EventType() string
	Version() int
	AggregateId() string
}

// IDomainEventHolder represents an entity that can hold and clear domain events.
type IDomainEventHolder[Event IDomainEvent] interface {
	AddDomainEvent(event Event)
	ClearDomainEvents()
	DomainEvents() []Event
}
