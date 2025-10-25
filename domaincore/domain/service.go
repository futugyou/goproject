package domain

import (
	"errors"
	"reflect"
)

type DomainService[E DomainEvent, ES EventSourcing] struct {
}

func NewDomainService[E DomainEvent, ES EventSourcing]() *DomainService[E, ES] {
	return &DomainService[E, ES]{}
}

func (ds *DomainService[E, ES]) RetrieveAllVersions(aggregate ES, events []E) ([]ES, error) {
	var aggregates []ES
	for _, event := range events {
		aggregate.Apply(event)

		aggregateValue := reflect.ValueOf(aggregate).Elem()
		tmp := reflect.New(aggregateValue.Type()).Interface().(ES)
		reflect.ValueOf(tmp).Elem().Set(aggregateValue)

		aggregates = append(aggregates, tmp)
	}
	return aggregates, nil
}

func (ds *DomainService[E, ES]) RetrieveSpecificVersion(aggregate ES, events []E, version int) (*ES, error) {
	if version < 0 {
		return nil, errors.New("invalid version number, must be non-negative")
	}

	if err := ds.applyEventsUntilVersion(aggregate, events, version); err != nil {
		return nil, err
	}

	if aggregate.AggregateVersion() != version {
		return nil, errors.New("the requested version is not available")
	}

	return &aggregate, nil
}

func (ds *DomainService[E, ES]) applyEventsUntilVersion(
	aggregate ES,
	events []E,
	targetVersion int,
) error {
	for _, event := range events {
		// Event version exceeds target version, stop immediately
		if event.Version() > targetVersion {
			break
		}

		// Event doesn't need to be applied, skip
		if !ds.shouldApplyEvent(event, aggregate, targetVersion) {
			continue
		}

		// Apply event
		aggregate.Apply(event)

		// Target version reached, stop loop
		if aggregate.AggregateVersion() == targetVersion {
			break
		}
	}
	return nil
}

func (ds *DomainService[E, ES]) shouldApplyEvent(event E, aggregate ES, targetVersion int) bool {
	eventVersion := event.Version()
	currentVersion := aggregate.AggregateVersion()
	return eventVersion > currentVersion && eventVersion <= targetVersion
}

func (ds *DomainService[E, ES]) RetrieveLatestVersion(aggregate ES, events []E) (*ES, error) {
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
