package infrastructure

import (
	"github.com/futugyou/infr-project/domain"
)

type IEventStore[Event domain.IDomainEvent] interface {
	Save(events []Event) error
	Load(id string) ([]Event, error)
}
