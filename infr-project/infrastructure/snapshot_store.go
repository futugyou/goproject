package infrastructure

import (
	"fmt"

	"github.com/futugyou/infr-project/domain"
)

type ISnapshotStore[EventSourcing domain.IEventSourcing] interface {
	LoadSnapshot(id string) (*EventSourcing, error)
	LoadSnapshotByVersion(id string, version int) (*EventSourcing, error)
	SaveSnapshot(aggregate EventSourcing) error
}

type MemorySnapshotStore[EventSourcing domain.IEventSourcing] struct {
	storage map[string][]EventSourcing
}

func NewMemorySnapshotStore[EventSourcing domain.IEventSourcing]() *MemorySnapshotStore[EventSourcing] {
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
