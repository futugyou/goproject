package application

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/futugyou/domaincore/domain"
	infra "github.com/futugyou/domaincore/infrastructure"
)

type AppService struct {
	unitOfWork domain.UnitOfWork
}

func NewAppService(
	unitOfWork domain.UnitOfWork,
) *AppService {
	return &AppService{
		unitOfWork: unitOfWork,
	}
}

func (s *AppService) WithUnitOfWork(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, err := s.unitOfWork.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	var rollbackErr error
	var commitErr error
	defer func() {
		if err != nil {
			rollbackErr = s.unitOfWork.Rollback(ctx)
		} else {
			commitErr = s.unitOfWork.Commit(ctx)
		}
	}()

	err = fn(ctx)
	return errors.Join(err, rollbackErr, commitErr)
}

type ApplicationService[Event domain.DomainEvent, EventSourcing domain.EventSourcing] struct {
	eventStore        infra.EventStore[Event]
	snapshotStore     infra.SnapshotStore[EventSourcing]
	innerService      *AppService
	domainService     *domain.DomainService[Event, EventSourcing]
	eventDispatcher    infra.EventDispatcher
	newAggregateFunc  func() EventSourcing
	needStoreSnapshot func(EventSourcing) bool
}

func NewApplicationService[Event domain.DomainEvent, EventSourcing domain.EventSourcing](
	eventStore infra.EventStore[Event],
	snapshotStore infra.SnapshotStore[EventSourcing],
	unitOfWork domain.UnitOfWork,
	newAggregateFunc func() EventSourcing,
	needStoreSnapshot func(EventSourcing) bool,
	eventDispatcher infra.EventDispatcher,
) *ApplicationService[Event, EventSourcing] {
	return &ApplicationService[Event, EventSourcing]{
		eventStore:        eventStore,
		snapshotStore:     snapshotStore,
		innerService:      NewAppService(unitOfWork),
		domainService:     domain.NewDomainService[Event, EventSourcing](),
		newAggregateFunc:  newAggregateFunc,
		needStoreSnapshot: needStoreSnapshot,
		eventDispatcher:    eventDispatcher,
	}
}

func (es *ApplicationService[Event, EventSourcing]) RetrieveAllVersions(ctx context.Context, id string) ([]EventSourcing, error) {
	events, err := es.eventStore.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("no event found with id: %s", id)
	}

	agg := es.newAggregateFunc()
	aggregate := &agg

	return es.domainService.RetrieveAllVersions(*aggregate, events)
}

func (es *ApplicationService[Event, EventSourcing]) RetrieveSpecificVersion(ctx context.Context, id string, version int) (*EventSourcing, error) {
	if version < 0 {
		return nil, errors.New("invalid version number, must be non-negative")
	}

	aggregate, err := es.RestoreFromSnapshotByVersion(ctx, id, version)
	if err != nil {
		return nil, err
	}

	if (*aggregate).AggregateVersion() < version {
		events, err := es.eventStore.LoadGreaterthanVersion(ctx, id, version)
		if err != nil {
			return nil, err
		}

		if len(events) == 0 {
			return nil, fmt.Errorf("can not found events with id: %s and version: %d", id, version)
		}

		return es.domainService.RetrieveSpecificVersion(*aggregate, events, version)
	}

	return aggregate, nil
}

func (es *ApplicationService[Event, EventSourcing]) RetrieveLatestVersion(ctx context.Context, id string) (*EventSourcing, error) {
	// Attempt to restore the latest snapshot
	aggregate, err := es.RestoreFromSnapshot(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load all events for the aggregate
	events, err := es.eventStore.LoadGreaterthanVersion(ctx, id, (*aggregate).AggregateVersion())
	if err != nil {
		return nil, err
	}

	return es.domainService.RetrieveLatestVersion(*aggregate, events)
}

func (s *ApplicationService[Event, EventSourcing]) SaveSnapshotAndEvent(ctx context.Context, aggregate EventSourcing) error {
	return s.SaveSnapshotAndMutileEvents(ctx, []EventSourcing{aggregate})
}

// this method can handle mutle version change
// PS: According to the current design, one operation will not generate multiple events, so the logic(foreach) of this method is unnecessary.
func (s *ApplicationService[Event, EventSourcing]) SaveSnapshotAndMutileEvents(ctx context.Context, aggregates []EventSourcing) error {
	es := aggregates[len(aggregates)-1].DomainEvents()
	events := make([]Event, 0)
	for i := range es {
		ev, ok := es[i].(Event)
		if !ok {
			return fmt.Errorf("expected type Event but got %T", es[i])
		}
		events = append(events, ev)
	}

	for _, agg := range aggregates {
		if s.needStoreSnapshot(agg) {
			if err := s.snapshotStore.SaveSnapshot(ctx, agg); err != nil {
				return err
			}
		}
	}

	s.eventDispatcher.DispatchDomainEvents(ctx, es)

	return s.eventStore.Save(ctx, events)
}

func (es *ApplicationService[Event, EventSourcing]) RestoreFromSnapshot(ctx context.Context, id string) (*EventSourcing, error) {
	datas, err := es.snapshotStore.LoadSnapshot(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(datas) == 0 {
		return nil, fmt.Errorf("can not found snapshot with id: %s", id)
	}

	return &datas[len(datas)-1], nil
}

func (es *ApplicationService[Event, EventSourcing]) RestoreFromSnapshotByVersion(ctx context.Context, id string, version int) (*EventSourcing, error) {
	if version < 0 {
		return nil, errors.New("invalid version number, must be non-negative")
	}

	datas, err := es.snapshotStore.LoadSnapshot(ctx, id)
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

func (s *ApplicationService[Event, EventSourcing]) WithUnitOfWork(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.innerService.WithUnitOfWork(ctx, fn)
}
