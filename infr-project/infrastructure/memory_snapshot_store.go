package infrastructure

import (
	"fmt"

	"github.com/futugyou/infr-project/domain"
)

type MemorySnapshotStore[EventSourcing domain.IEventSourcing] struct {
	storage map[string][]EventSourcing
}

func NewMemorySnapshotStore[EventSourcing domain.IEventSourcing]() *MemorySnapshotStore[EventSourcing] {
	return &MemorySnapshotStore[EventSourcing]{
		storage: make(map[string][]EventSourcing),
	}
}

func (s *MemorySnapshotStore[EventSourcing]) LoadSnapshot(id string) ([]EventSourcing, error) {
	datas, ok := s.storage[id]
	if !ok || len(datas) == 0 {
		return nil, fmt.Errorf("no data for %s", id)
	}

	return datas, nil
}

func (s *MemorySnapshotStore[EventSourcing]) SaveSnapshot(aggregate EventSourcing) error {
	s.storage[aggregate.AggregateId()] = append(s.storage[aggregate.AggregateId()], aggregate)
	return nil
}
