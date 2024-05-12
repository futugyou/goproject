package domain

import (
	"errors"
)

type GeneralEventSourcer[Event IDomainEvent, EventSourcing IEventSourcing] struct {
	IEventStore[Event]
	ISnapshotStore[EventSourcing]
}

func NewEventSourcer[Event IDomainEvent, EventSourcing IEventSourcing]() *GeneralEventSourcer[Event, EventSourcing] {
	return &GeneralEventSourcer[Event, EventSourcing]{
		IEventStore:    NewMemoryEventStore[Event](),
		ISnapshotStore: NewMemorySnapshotStore[EventSourcing](),
	}
}

func (es *GeneralEventSourcer[Event, EventSourcing]) Apply(aggregate EventSourcing, event Event) error {
	return aggregate.Apply(event)
}

func (es *GeneralEventSourcer[Event, EventSourcing]) GetAllVersions(id string) ([]EventSourcing, error) {
	events, err := es.Load(id)
	if err != nil {
		return nil, err
	}

	var aggregates []EventSourcing
	aggregate := *new(EventSourcing)
	for _, event := range events {
		es.Apply(aggregate, event)
		aggregates = append(aggregates, aggregate)

	}
	return aggregates, nil
}

func (es *GeneralEventSourcer[Event, EventSourcing]) GetSpecificVersion(id string, version int) (*EventSourcing, error) {
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
			aggregate = new(EventSourcing)
		}

		// Apply events to the snapshot or the newly created aggregate until the requested version is reached
		for _, event := range events {
			eventVersion := event.Version()
			if eventVersion > (*aggregate).AggregateVersion() && eventVersion <= version {
				es.Apply(*aggregate, event)
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

func (es *GeneralEventSourcer[Event, EventSourcing]) GetLatestVersion(id string) (*EventSourcing, error) {
	// Attempt to restore the latest snapshot
	aggregate, err := es.RestoreFromSnapshot(id)
	if err != nil {
		// If an error occurs, we assume no snapshot is available and start from scratch
		aggregate = new(EventSourcing)
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
			es.Apply(*aggregate, event)
		}
	}

	// The aggregate is now at the latest version
	return aggregate, nil
}

func (es *GeneralEventSourcer[Event, EventSourcing]) TakeSnapshot(aggregate EventSourcing) error {
	return es.SaveSnapshot(aggregate)
}

func (es *GeneralEventSourcer[Event, EventSourcing]) RestoreFromSnapshot(id string) (*EventSourcing, error) {
	return es.LoadSnapshot(id)
}

func (es *GeneralEventSourcer[Event, EventSourcing]) RestoreFromSnapshotByVersion(id string, version int) (*EventSourcing, error) {
	return es.LoadSnapshotByVersion(id, version)
}
