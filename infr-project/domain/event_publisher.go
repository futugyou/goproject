package domain

import (
	"context"
)

type IEventPublisher interface {
	Publish(ctx context.Context, events []IDomainEvent) error
	PublishCommon(ctx context.Context, event any, event_type string) error
}
