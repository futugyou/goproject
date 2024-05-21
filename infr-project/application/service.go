package application

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/futugyou/infr-project/domain"
	infra "github.com/futugyou/infr-project/infrastructure"
)

type ApplicationService[Event domain.IDomainEvent, EventSourcing domain.IEventSourcing] struct {
	eventStore        infra.IEventStore[Event]
	snapshotStore     infra.ISnapshotStore[EventSourcing]
	unitOfWork        domain.IUnitOfWork
	domainService     *domain.DomainService[Event, EventSourcing]
	newAggregateFunc  func() EventSourcing
	needStoreSnapshot func(EventSourcing) bool
}

func NewApplicationService[Event domain.IDomainEvent, EventSourcing domain.IEventSourcing](
	eventStore infra.IEventStore[Event],
	snapshotStore infra.ISnapshotStore[EventSourcing],
	unitOfWork domain.IUnitOfWork,
	newAggregateFunc func() EventSourcing,
	needStoreSnapshot func(EventSourcing) bool,
) *ApplicationService[Event, EventSourcing] {
	return &ApplicationService[Event, EventSourcing]{
		eventStore:        eventStore,
		snapshotStore:     snapshotStore,
		unitOfWork:        unitOfWork,
		domainService:     domain.NewDomainService[Event, EventSourcing](),
		newAggregateFunc:  newAggregateFunc,
		needStoreSnapshot: needStoreSnapshot,
	}
}

func (es *ApplicationService[Event, EventSourcing]) RetrieveAllVersions(id string) ([]EventSourcing, error) {
	events, err := es.eventStore.Load(id)
	if err != nil {
		return nil, err
	}

	agg := es.newAggregateFunc()
	aggregate := &agg

	return es.domainService.RetrieveAllVersions(*aggregate, events)
}

func (es *ApplicationService[Event, EventSourcing]) RetrieveSpecificVersion(id string, version int) (*EventSourcing, error) {
	if version < 0 {
		return nil, errors.New("invalid version number, must be non-negative")
	}

	aggregate, err := es.RestoreFromSnapshotByVersion(id, version)
	if err != nil || (*aggregate).AggregateVersion() < version {
		events, err := es.eventStore.Load(id)
		if err != nil {
			return nil, err
		}

		if aggregate == nil {
			agg := es.newAggregateFunc()
			aggregate = &agg
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
		agg := es.newAggregateFunc()
		aggregate = &agg
	}

	// Load all events for the aggregate
	events, err := es.eventStore.Load(id)
	if err != nil {
		return nil, err
	}

	return es.domainService.RetrieveLatestVersion(*aggregate, events)
}

func (s *ApplicationService[Event, EventSourcing]) SaveSnapshotAndEvent(ctx context.Context, aggregate EventSourcing) error {
	es := aggregate.DomainEvents()
	events := make([]Event, 0)
	for i := 0; i < len(es); i++ {
		ev, ok := es[i].(Event)
		if !ok {
			return fmt.Errorf("expected type Event but got %T", es[i])
		}
		events = append(events, ev)
	}

	if err := s.snapshotStore.SaveSnapshot(ctx, aggregate); err != nil {
		return err
	}

	return s.eventStore.Save(ctx, events)
}

// Deprecated: TakeSnapshot is deprecated. it is not necessary to use it alone, may be use SaveSnapshotAndEvent is a good idea.
func (es *ApplicationService[Event, EventSourcing]) TakeSnapshot(aggregate EventSourcing) error {
	// aggregate is created with version 1
	// The current storage snapshot logic starts from the first version and is saved every 5 versions.
	if aggregate.AggregateVersion()%5 != 1 {
		return nil
	}
	return es.snapshotStore.SaveSnapshot(context.Background(), aggregate)
}

func (es *ApplicationService[Event, EventSourcing]) RestoreFromSnapshot(id string) (*EventSourcing, error) {
	datas, err := es.snapshotStore.LoadSnapshot(id)
	if err != nil {
		return nil, err
	}

	if len(datas) == 0 {
		return nil, fmt.Errorf("can not found snapshot with id: %s", id)
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

func (s *ApplicationService[Event, EventSourcing]) withUnitOfWork(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, err := s.unitOfWork.Start(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			s.unitOfWork.Rollback(ctx)
		} else {
			commitErr := s.unitOfWork.Commit(ctx)
			if commitErr != nil {
				err = commitErr
			}
		}
		s.unitOfWork.End(ctx)
	}()

	err = fn(ctx)
	return err
}
