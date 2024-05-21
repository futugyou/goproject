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

type DomainEvent struct {
	Id              string `bson:"id" json:"id"`
	ResourceVersion int    `bson:"version" json:"version"`
	EventType       string `bson:"event_type" json:"event_type"`
}

func (d DomainEvent) Version() int {
	return d.ResourceVersion
}

func (d DomainEvent) AggregateId() string {
	return d.Id
}
