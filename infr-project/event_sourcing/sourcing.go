package eventsourcing

import (
	"errors"
)

// IEvent represents the interface for events.
type IEvent interface {
	EventType() string
}

// IAggregate represents the basic interface for aggregates.
type IAggregate interface {
	AggregateName() string
	AggregateId() string
}

// IEventSourcing extends IAggregate with event sourcing specific methods, including versioning.
type IEventSourcing interface {
	IAggregate
	Apply(event IEvent) error
	AggregateVersion() int
}

type IEventSourcer[E IEvent, R IEventSourcing] interface {
	Save(events []E) error
	Load(id string) ([]E, error)
	Apply(aggregate R, event E) R
	GetAllVersions(id string) ([]R, error)
	GetSpecificVersion(id string, version int) (*R, error)
}

type GeneralEventSourcer[E IEvent, R IEventSourcing] struct {
	storage       IEventStorage[E]
	snapshotStore ISnapshotStore[R]
}

func NewEventSourcer[E IEvent, R IEventSourcing]() *GeneralEventSourcer[E, R] {
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

	// TODO (*aggregate).AggregateVersion() < version

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
