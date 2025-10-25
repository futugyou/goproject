package domain

import "time"

// DomainEvent represents the interface for events.
type DomainEvent interface {
	EventType() string
	Version() int
	AggregateID() string
}

// DomainEventHolder represents an entity that can hold and clear domain events.
type DomainEventHolder[Event DomainEvent] interface {
	AggregateRoot
	AddDomainEvent(event Event)
	ClearDomainEvents()
	DomainEvents() []Event
}

type BaseDomainEvent struct {
	ID              string
	ResourceVersion int
	EventType       string
	CreatedAt       time.Time
}

func (d BaseDomainEvent) Version() int {
	return d.ResourceVersion
}

func (d BaseDomainEvent) AggregateID() string {
	return d.ID
}
