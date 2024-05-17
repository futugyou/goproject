package domain

import "errors"

type DomainService[Event IDomainEvent, EventSourcing IEventSourcing] struct {
}

func NewDomainService[Event IDomainEvent, EventSourcing IEventSourcing]() *DomainService[Event, EventSourcing] {
	return &DomainService[Event, EventSourcing]{}
}

func (ds *DomainService[Event, EventSourcing]) RetrieveAllVersions(aggregate EventSourcing, events []Event) ([]EventSourcing, error) {
	var aggregates []EventSourcing
	for _, event := range events {
		aggregate.Apply(event)
		aggregates = append(aggregates, aggregate)
	}
	return aggregates, nil
}

func (ds *DomainService[Event, EventSourcing]) RetrieveSpecificVersion(aggregate EventSourcing, events []Event, version int) (*EventSourcing, error) {
	if version < 0 {
		return nil, errors.New("invalid version number, must be non-negative")
	}

	for _, event := range events {
		eventVersion := event.Version()
		if eventVersion > aggregate.AggregateVersion() && eventVersion <= version {
			aggregate.Apply(event)
			if aggregate.AggregateVersion() == version {
				break
			}
		} else if eventVersion > version {
			break // Stop processing if the event's version surpasses the target version
		}
	}

	// After applying events, check if the current version of the aggregate matches the requested version
	if aggregate.AggregateVersion() != version {
		return nil, errors.New("the requested version is not available")
	}

	return &aggregate, nil
}

func (ds *DomainService[Event, EventSourcing]) RetrieveLatestVersion(aggregate EventSourcing, events []Event) (*EventSourcing, error) {
	// Apply events to the snapshot or the newly created aggregate
	for _, event := range events {
		eventVersion := event.Version()
		// Only apply events that are newer than the snapshot's version
		if eventVersion > aggregate.AggregateVersion() {
			aggregate.Apply(event)
		}
	}

	// The aggregate is now at the latest version
	return &aggregate, nil
}
