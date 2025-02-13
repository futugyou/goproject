package infrastructure

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type IEventPublisher interface {
	Publish(ctx context.Context, events []domain.IDomainEvent) error
	PublishCommon(ctx context.Context, event any, event_type string) error
}
