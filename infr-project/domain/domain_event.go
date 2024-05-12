package domain

// IEvent represents the interface for events.
type IDomainEvent interface {
	EventType() string
	Version() int
}

// IDomainEventHolder represents an entity that can hold and clear domain events.
type IDomainEventHolder interface {
	AddDomainEvent(event IDomainEvent)
	ClearDomainEvents()
	DomainEvents() []IDomainEvent
}
