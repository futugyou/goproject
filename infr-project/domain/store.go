package domain

import "fmt"

type IEventStore[Event IDomainEvent] interface {
	Save(events []Event) error
	Load(id string) ([]Event, error)
}

type MemoryEventStore[Event IDomainEvent] struct {
	storage map[string][]Event
}

func NewMemoryEventStore[Event IDomainEvent]() *MemoryEventStore[Event] {
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

type ISnapshotStore[EventSourcing IEventSourcing] interface {
	LoadSnapshot(id string) (*EventSourcing, error)
	LoadSnapshotByVersion(id string, version int) (*EventSourcing, error)
	SaveSnapshot(aggregate EventSourcing) error
}

type MemorySnapshotStore[EventSourcing IEventSourcing] struct {
	storage map[string][]EventSourcing
}

func NewMemorySnapshotStore[EventSourcing IEventSourcing]() *MemorySnapshotStore[EventSourcing] {
	return &MemorySnapshotStore[EventSourcing]{
		storage: make(map[string][]EventSourcing),
	}
}

func (s *MemorySnapshotStore[EventSourcing]) LoadSnapshot(id string) (*EventSourcing, error) {
	datas, ok := s.storage[id]
	if !ok || len(datas) == 0 {
		return nil, fmt.Errorf("no data for %s", id)
	}

	return &datas[len(datas)-1], nil
}

func (s *MemorySnapshotStore[R]) LoadSnapshotByVersion(id string, version int) (*R, error) {
	datas, ok := s.storage[id]
	if !ok {
		return nil, fmt.Errorf("no data for %s", id)
	}
	for i := len(datas) - 1; i >= 0; i-- {
		if datas[i].AggregateVersion() <= version { //TODO it will change to '<' after add event verion
			return &datas[i], nil
		}
	}
	return nil, fmt.Errorf("no data for id %s version %d", id, version)
}

func (s *MemorySnapshotStore[EventSourcing]) SaveSnapshot(aggregate EventSourcing) error {
	if aggregate.AggregateVersion()%5 != 0 {
		return nil
	}
	s.storage[aggregate.AggregateId()] = append(s.storage[aggregate.AggregateId()], aggregate)
	return nil
}
