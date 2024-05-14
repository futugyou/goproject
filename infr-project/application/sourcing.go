package application

import (
	domain "github.com/futugyou/infr-project/domain"
	infr "github.com/futugyou/infr-project/infrastructure"
)

type IEventSourcingService[Event domain.IDomainEvent, EventSourcing domain.IEventSourcing] interface {
	infr.IEventStore[Event]
	domain.IEventApplier[Event, EventSourcing]
	domain.IAggregateRetriever[EventSourcing]
}
