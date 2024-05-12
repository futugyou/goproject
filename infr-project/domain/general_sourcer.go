package domain

import (
	"errors"
	"fmt"
)

type GeneralEventSourcer[E IDomainEvent, R IEventSourcing] struct {
	IEventStore[E]
	ISnapshotStore[R]
}

func NewEventSourcer[E IDomainEvent, R IEventSourcing]() *GeneralEventSourcer[E, R] {
	return &GeneralEventSourcer[E, R]{
		IEventStore:    NewMemoryEventStore[E](),
		ISnapshotStore: NewMemorySnapshotStore[R](),
	}
}

func (es *GeneralEventSourcer[E, R]) Apply(aggregate R, event E) (R, error) {
	newAggregate, err := aggregate.Apply(event)
	if err != nil {
		return aggregate, err
	}

	typedAggregate, ok := newAggregate.(R)
	if !ok {
		return aggregate, fmt.Errorf("apply method returned an instance of incorrect type")
	}

	return typedAggregate, nil
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
		return nil, errors.New("invalid version number, must be non-negative")
	}

	aggregate, err := es.RestoreFromSnapshotByVersion(id, version)
	if err != nil || (*aggregate).AggregateVersion() < version {
		events, eventsErr := es.Load(id)
		if eventsErr != nil {
			return nil, eventsErr
		}

		if aggregate == nil {
			// Initialize an empty aggregate if snapshot doesn't exist or an error occurred
			aggregate = new(R)
		}

		// Apply events to the snapshot or the newly created aggregate until the requested version is reached
		for _, event := range events {
			eventVersion := event.Version()
			if eventVersion > (*aggregate).AggregateVersion() && eventVersion <= version {
				(*aggregate).Apply(event)
				if (*aggregate).AggregateVersion() == version {
					break
				}
			} else if eventVersion > version {
				break // Stop processing if the event's version surpasses the target version
			}
		}

		// After applying events, check if the current version of the aggregate matches the requested version
		if (*aggregate).AggregateVersion() != version {
			return nil, errors.New("the requested version is not available")
		}
	}

	return aggregate, nil
}

func (es *GeneralEventSourcer[E, R]) GetLatestVersion(id string) (*R, error) {
	// Attempt to restore the latest snapshot
	aggregate, err := es.RestoreFromSnapshot(id)
	if err != nil {
		// If an error occurs, we assume no snapshot is available and start from scratch
		aggregate = new(R)
	}

	// Load all events for the aggregate
	events, eventsErr := es.Load(id)
	if eventsErr != nil {
		return nil, eventsErr
	}

	// Apply events to the snapshot or the newly created aggregate
	for _, event := range events {
		eventVersion := event.Version()
		// Only apply events that are newer than the snapshot's version
		if eventVersion > (*aggregate).AggregateVersion() {
			(*aggregate).Apply(event)
		}
	}

	// The aggregate is now at the latest version
	return aggregate, nil
}

func (es *GeneralEventSourcer[E, R]) TakeSnapshot(aggregate R) error {
	return es.SaveSnapshot(aggregate)
}

func (es *GeneralEventSourcer[E, R]) RestoreFromSnapshot(id string) (*R, error) {
	return es.LoadSnapshot(id)
}

func (es *GeneralEventSourcer[E, R]) RestoreFromSnapshotByVersion(id string, version int) (*R, error) {
	return es.LoadSnapshotByVersion(id, version)
}
