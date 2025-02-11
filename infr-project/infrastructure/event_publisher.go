package infrastructure

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type IEventPulisher[Event domain.IDomainEvent] interface {
	Publish(ctx context.Context, events []Event) error
}
