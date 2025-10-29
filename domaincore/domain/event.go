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
	ID              string    `bson:"id" redis:"id" json:"id"`
	ResourceVersion int       `bson:"version" redis:"version" json:"version"`
	EventType       string    `bson:"event_type" redis:"event_type" json:"event_type"`
	CreatedAt       time.Time `bson:"created_at" redis:"created_at" json:"created_at"`
}

func (d BaseDomainEvent) Version() int {
	return d.ResourceVersion
}

func (d BaseDomainEvent) AggregateID() string {
	return d.ID
}
