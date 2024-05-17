package application

import (
	domain "github.com/futugyou/infr-project/domain"
	infra "github.com/futugyou/infr-project/infrastructure"
)

// Deprecated: IEventSourcingService is deprecated.
type IEventSourcingService[Event domain.IDomainEvent, EventSourcing domain.IEventSourcing] interface {
	infra.IEventStore[Event]
	domain.ISnapshotter[Event, EventSourcing]
	domain.IAggregateRetriever[EventSourcing]
}
