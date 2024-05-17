package application

import (
	"errors"
	"fmt"

	domain "github.com/futugyou/infr-project/domain"
	infra "github.com/futugyou/infr-project/infrastructure"
)

type ApplicationService[Event domain.IDomainEvent, EventSourcing domain.IEventSourcing] struct {
	eventStore    infra.IEventStore[Event]
	snapshotStore infra.ISnapshotStore[EventSourcing]
	domainService *domain.DomainService[Event, EventSourcing]
	instance      EventSourcing
}

func NewApplicationService[Event domain.IDomainEvent, EventSourcing domain.IEventSourcing](
	eventStore infra.IEventStore[Event],
	snapshotStore infra.ISnapshotStore[EventSourcing],
	instance EventSourcing,
) *ApplicationService[Event, EventSourcing] {
	return &ApplicationService[Event, EventSourcing]{
		eventStore:    eventStore,
		snapshotStore: snapshotStore,
		domainService: domain.NewDomainService[Event, EventSourcing](),
		instance:      instance,
	}
}

func (es *ApplicationService[Event, EventSourcing]) RetrieveAllVersions(id string) ([]EventSourcing, error) {
	events, err := es.eventStore.Load(id)
	if err != nil {
		return nil, err
	}
	aggregate, err := es.RestoreFromSnapshot(id)
	if err != nil || aggregate == nil {
		// If no snapshot is available or an error occurs, start from scratch
		aggregate = &es.instance
	}
	return es.domainService.RetrieveAllVersions(*aggregate, events)
}

func (es *ApplicationService[Event, EventSourcing]) RetrieveSpecificVersion(id string, version int) (*EventSourcing, error) {
	if version < 0 {
		return nil, errors.New("invalid version number, must be non-negative")
	}

	aggregate, err := es.RestoreFromSnapshotByVersion(id, version)
	if err != nil || (*aggregate).AggregateVersion() < version {
		events, eventsErr := es.eventStore.Load(id)
		if eventsErr != nil {
			return nil, eventsErr
		}

		if aggregate == nil {
			// Initialize an empty aggregate if snapshot doesn't exist or an error occurred
			aggregate = &es.instance
		}

		return es.domainService.RetrieveSpecificVersion(*aggregate, events, version)
	}

	return aggregate, nil
}

func (es *ApplicationService[Event, EventSourcing]) RetrieveLatestVersion(id string) (*EventSourcing, error) {
	// Attempt to restore the latest snapshot
	aggregate, err := es.RestoreFromSnapshot(id)
	if err != nil || aggregate == nil {
		// If an error occurs, we assume no snapshot is available and start from scratch
		aggregate = &es.instance
	}

	// Load all events for the aggregate
	events, err := es.eventStore.Load(id)
	if err != nil {
		return nil, err
	}

	return es.domainService.RetrieveLatestVersion(*aggregate, events)
}

func (es *ApplicationService[Event, EventSourcing]) TakeSnapshot(aggregate EventSourcing) error {
	// aggregate is created with version 1
	// The current storage snapshot logic starts from the first version and is saved every 5 versions.
	if aggregate.AggregateVersion()%5 != 1 {
		return nil
	}
	return es.snapshotStore.SaveSnapshot(aggregate)
}

func (es *ApplicationService[Event, EventSourcing]) RestoreFromSnapshot(id string) (*EventSourcing, error) {
	datas, err := es.snapshotStore.LoadSnapshot(id)
	if err != nil {
		return nil, err
	}

	return &datas[len(datas)-1], nil
}

func (es *ApplicationService[Event, EventSourcing]) RestoreFromSnapshotByVersion(id string, version int) (*EventSourcing, error) {
	if version < 0 {
		return nil, errors.New("invalid version number, must be non-negative")
	}

	datas, err := es.snapshotStore.LoadSnapshot(id)
	if err != nil {
		return nil, err
	}

	for i := len(datas); i >= 0; i-- {
		if datas[i].AggregateVersion() <= version {
			return &datas[i], nil
		}
	}

	return nil, fmt.Errorf("can not found snapshot with id: %s and version: %d", id, version)
}
