package infrastructure

import (
	"context"

	"github.com/futugyou/domaincore/domain"
)

type Event interface {
	EventType() string
}

type EventDispatcher interface {
	DispatchDomainEvents(ctx context.Context, events []domain.DomainEvent) error
	DispatchIntegrationEvent(ctx context.Context, event Event) error
}
