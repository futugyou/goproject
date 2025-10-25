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
