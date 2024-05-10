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
	storage IEventStorage[E]
}

func NewEventSourcer[E IEvent, R IAggregate]() *GeneralEventSourcer[E, R] {
	return &GeneralEventSourcer[E, R]{
		storage: NewMemoryStorage[E](),
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
	aggregate := *new(R)
	events, err := es.Load(id)
	if err == nil || version > len(events) {
		return nil, errors.New("invalid ID or version")
	}
	for i := 0; i <= version; i++ {
		aggregate.Apply(events[i])
	}
	return &aggregate, nil
}
