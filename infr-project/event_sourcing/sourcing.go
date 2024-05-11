package eventsourcing

import (
	"errors"
)

type IEvent interface {
	EventType() string
}

type IAggregate interface {
	AggregateName() string
	AggregateId() string
	AggregateVersion() int
	Apply(event IEvent) error
}

type IEventSourcer[E IEvent, R IAggregate] interface {
	Save(events []E) error
	Load(id string) ([]E, error)
	Apply(aggregate R, event E) R
	GetAllVersions(id string) ([]R, error)
	GetSpecificVersion(id string, version int) (*R, error)
}

type GeneralEventSourcer[E IEvent, R IAggregate] struct {
	storage       IEventStorage[E]
	snapshotStore ISnapshotStore[R]
}

func NewEventSourcer[E IEvent, R IAggregate]() *GeneralEventSourcer[E, R] {
	return &GeneralEventSourcer[E, R]{
		storage:       NewMemoryStorage[E](),
		snapshotStore: NewMemorySnapshotStore[R](),
	}
}

func (es *GeneralEventSourcer[E, R]) Save(events []E) error {
	return es.storage.SaveEvents(events)
}

func (es *GeneralEventSourcer[E, R]) Load(id string) ([]E, error) {
	events, err := es.storage.GetEvents(id)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (es *GeneralEventSourcer[E, R]) Apply(aggregate R, event E) R {
	aggregate.Apply(event)
	return aggregate
}

func (es *GeneralEventSourcer[E, R]) GetAllVersions(id string) ([]R, error) {
	events, err := es.Load(id)
	if err != nil {
		return nil, err
	}

	var aggregates []R
	aggregate := *new(R)
	for _, event := range events {
		aggregate.Apply(event)
		aggregates = append(aggregates, aggregate)

	}
	return aggregates, nil
}

func (es *GeneralEventSourcer[E, R]) GetSpecificVersion(id string, version int) (*R, error) {
	if version < 0 {
		return nil, errors.New("invalid ID or version")
	}
	aggregate, err := es.RestoreFromSnapshotByVersion(id, version)
	if err != nil {
		events, err := es.Load(id)
		if err == nil || version > len(events) {
			return nil, errors.New("invalid ID or version")
		}
		for i := 1; i <= version; i++ {
			(*aggregate).Apply(events[i])
		}
	}
	// TODO aggregate.version < version

	return aggregate, nil
}

func (es *GeneralEventSourcer[E, R]) TakeSnapshot(aggregate R) error {
	return es.snapshotStore.SaveSnapshot(aggregate)
}

func (es *GeneralEventSourcer[E, R]) RestoreFromSnapshot(id string) (*R, error) {
	return es.snapshotStore.LoadSnapshot(id)
}

func (es *GeneralEventSourcer[E, R]) RestoreFromSnapshotByVersion(id string, version int) (*R, error) {
	return es.snapshotStore.LoadSnapshotByVersion(id, version)
}
