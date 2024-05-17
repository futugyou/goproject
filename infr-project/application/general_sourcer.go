package application

import (
	"errors"

	domain "github.com/futugyou/infr-project/domain"
	infra "github.com/futugyou/infr-project/infrastructure"
)

// Deprecated: GeneralEventSourcer is deprecated, use ApplicationService instead.
type GeneralEventSourcer[Event domain.IDomainEvent, EventSourcing domain.IEventSourcing] struct {
	infra.IEventStore[Event]
	infra.ISnapshotStore[EventSourcing]
	instance EventSourcing
}

// Deprecated: GeneralEventSourcer is deprecated, so NewEventSourcer is deprecated, use ApplicationService instead.
func NewEventSourcer[Event domain.IDomainEvent, EventSourcing domain.IEventSourcing](
	eventStore infra.IEventStore[Event],
	snapshotStore infra.ISnapshotStore[EventSourcing],
	instance EventSourcing,
) *GeneralEventSourcer[Event, EventSourcing] {
	return &GeneralEventSourcer[Event, EventSourcing]{
		IEventStore:    eventStore,
		ISnapshotStore: snapshotStore,
		instance:       instance,
	}
}

func (es *GeneralEventSourcer[Event, EventSourcing]) RetrieveAllVersions(id string) ([]EventSourcing, error) {
	events, err := es.Load(id)
	if err != nil {
		return nil, err
	}
	var aggregates []EventSourcing
	aggregate := es.instance
	for _, event := range events {
		aggregate.Apply(event)
		aggregates = append(aggregates, aggregate)

	}
	return aggregates, nil
}

func (es *GeneralEventSourcer[Event, EventSourcing]) RetrieveSpecificVersion(id string, version int) (*EventSourcing, error) {
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
			aggregate = &es.instance
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

func (es *GeneralEventSourcer[Event, EventSourcing]) RetrieveLatestVersion(id string) (*EventSourcing, error) {
	// Attempt to restore the latest snapshot
	aggregate, err := es.RestoreFromSnapshot(id)
	if err != nil || aggregate == nil {
		// If an error occurs, we assume no snapshot is available and start from scratch
		aggregate = &es.instance
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

func (es *GeneralEventSourcer[Event, EventSourcing]) TakeSnapshot(aggregate EventSourcing) error {
	return es.SaveSnapshot(aggregate)
}

func (es *GeneralEventSourcer[Event, EventSourcing]) RestoreFromSnapshot(id string) (*EventSourcing, error) {
	return es.LoadSnapshot(id)
}

func (es *GeneralEventSourcer[Event, EventSourcing]) RestoreFromSnapshotByVersion(id string, version int) (*EventSourcing, error) {
	return es.LoadSnapshotByVersion(id, version)
}
