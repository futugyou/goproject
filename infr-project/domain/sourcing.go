package domain

// IEventSourcing extends IAggregate with event sourcing specific methods, including versioning.
type IEventSourcing interface {
	IAggregate
	IDomainEventHolder[IDomainEvent]
	Apply(event IDomainEvent) error
	AggregateVersion() int
}

type BaseEventSourcing struct {
	domainEvents []IDomainEvent `json:"-"`
}

func (b *BaseEventSourcing) AddDomainEvent(event IDomainEvent) {
	b.domainEvents = append(b.domainEvents, event)
}

func (b *BaseEventSourcing) ClearDomainEvents() {
	b.domainEvents = nil
}

func (b *BaseEventSourcing) DomainEvents() []IDomainEvent {
	return b.domainEvents
}
