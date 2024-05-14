package infrastructure

import (
	"fmt"

	"github.com/futugyou/infr-project/domain"
)

type IEventStore[Event domain.IDomainEvent] interface {
	Save(events []Event) error
	Load(id string) ([]Event, error)
}

type MemoryEventStore[Event domain.IDomainEvent] struct {
	storage map[string][]Event
}

func NewMemoryEventStore[Event domain.IDomainEvent]() *MemoryEventStore[Event] {
	return &MemoryEventStore[Event]{
		storage: make(map[string][]Event),
	}
}

func (s *MemoryEventStore[Event]) Load(id string) ([]Event, error) {
	events, ok := s.storage[id]
	if !ok {
		return nil, fmt.Errorf("no data for %s", id)
	}

	return events, nil
}

func (s *MemoryEventStore[Event]) Save(events []Event) error {
	for _, event := range events {
		id := event.EventType()
		s.storage[id] = append(s.storage[id], event)
	}
	return nil
}
