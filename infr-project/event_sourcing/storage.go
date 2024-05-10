package eventsourcing

import "fmt"

type IEventStorage[E IEvent] interface {
	GetEvents(aggregateId string) ([]E, error)
	SaveEvents(events []E) error
}

type MemoryStorage[E IEvent] struct {
	storage map[string][]E
}

func NewMemoryStorage[E IEvent]() *MemoryStorage[E] {
	return &MemoryStorage[E]{
		storage: make(map[string][]E),
	}
}

func (s *MemoryStorage[E]) GetEvents(aggregateId string) ([]E, error) {
	events, ok := s.storage[aggregateId]
	if !ok {
		return nil, fmt.Errorf("no data for %s", aggregateId)
	}

	return events, nil
}

func (s *MemoryStorage[E]) SaveEvents(events []E) error {
	for _, event := range events {
		id := event.EventType()
		s.storage[id] = append(s.storage[id], event)
	}
	return nil
}
