package domain

import "fmt"

type IEventStore[E IDomainEvent] interface {
	Save(events []E) error
	Load(id string) ([]E, error)
}

type MemoryEventStore[E IDomainEvent] struct {
	storage map[string][]E
}

func NewMemoryEventStore[E IDomainEvent]() *MemoryEventStore[E] {
	return &MemoryEventStore[E]{
		storage: make(map[string][]E),
	}
}

func (s *MemoryEventStore[E]) Load(id string) ([]E, error) {
	events, ok := s.storage[id]
	if !ok {
		return nil, fmt.Errorf("no data for %s", id)
	}

	return events, nil
}

func (s *MemoryEventStore[E]) Save(events []E) error {
	for _, event := range events {
		id := event.EventType()
		s.storage[id] = append(s.storage[id], event)
	}
	return nil
}

type ISnapshotStore[R IEventSourcing] interface {
	LoadSnapshot(id string) (*R, error)
	LoadSnapshotByVersion(id string, version int) (*R, error)
	SaveSnapshot(aggregate R) error
}

type MemorySnapshotStore[R IEventSourcing] struct {
	storage map[string][]R
}

func NewMemorySnapshotStore[R IEventSourcing]() *MemorySnapshotStore[R] {
	return &MemorySnapshotStore[R]{
		storage: make(map[string][]R),
	}
}

func (s *MemorySnapshotStore[R]) LoadSnapshot(id string) (*R, error) {
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

func (s *MemorySnapshotStore[R]) SaveSnapshot(aggregate R) error {
	if aggregate.AggregateVersion()%5 != 0 {
		return nil
	}
	s.storage[aggregate.AggregateId()] = append(s.storage[aggregate.AggregateId()], aggregate)
	return nil
}
