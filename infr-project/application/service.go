package application

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/futugyou/infr-project/domain"
	infra "github.com/futugyou/infr-project/infrastructure"
)

type AppService struct {
	unitOfWork domain.IUnitOfWork
}

func NewAppService(
	unitOfWork domain.IUnitOfWork,
) *AppService {
	return &AppService{
		unitOfWork: unitOfWork,
	}
}

func (s *AppService) withUnitOfWork(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, err := s.unitOfWork.Start(ctx)
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

type ApplicationService[Event domain.IDomainEvent, EventSourcing domain.IEventSourcing] struct {
	eventStore        infra.IEventStore[Event]
	snapshotStore     infra.ISnapshotStore[EventSourcing]
	innerService      *AppService
	domainService     *domain.DomainService[Event, EventSourcing]
	eventPulisher     infra.IEventPulisher
	newAggregateFunc  func() EventSourcing
	needStoreSnapshot func(EventSourcing) bool
}

func NewApplicationService[Event domain.IDomainEvent, EventSourcing domain.IEventSourcing](
	eventStore infra.IEventStore[Event],
	snapshotStore infra.ISnapshotStore[EventSourcing],
	unitOfWork domain.IUnitOfWork,
	newAggregateFunc func() EventSourcing,
	needStoreSnapshot func(EventSourcing) bool,
	eventPulisher infra.IEventPulisher,
) *ApplicationService[Event, EventSourcing] {
	return &ApplicationService[Event, EventSourcing]{
		eventStore:        eventStore,
		snapshotStore:     snapshotStore,
		innerService:      NewAppService(unitOfWork),
		domainService:     domain.NewDomainService[Event, EventSourcing](),
		newAggregateFunc:  newAggregateFunc,
		needStoreSnapshot: needStoreSnapshot,
		eventPulisher:     eventPulisher,
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

// this method can not handle mutile version change
// eg. if version changed form 4 to 7, than event 4567 will be stored, but the snapshot(version 6) will not stored,
// beacuse the current aggregate version is 7, and needStoreSnapshot will return false(7%5==1)
func (s *ApplicationService[Event, EventSourcing]) SaveSnapshotAndEvent(ctx context.Context, aggregate EventSourcing) error {
	// es := aggregate.DomainEvents()
	// events := make([]Event, 0)
	// for i := 0; i < len(es); i++ {
	// 	ev, ok := es[i].(Event)
	// 	if !ok {
	// 		return fmt.Errorf("expected type Event but got %T", es[i])
	// 	}
	// 	events = append(events, ev)
	// }

	// if s.needStoreSnapshot(aggregate) {
	// 	if err := s.snapshotStore.SaveSnapshot(ctx, aggregate); err != nil {
	// 		return err
	// 	}
	// }

	// return s.eventStore.Save(ctx, events)
	return s.SaveSnapshotAndEvent2(ctx, []EventSourcing{aggregate})
}

// this method can handle mutle version change
func (s *ApplicationService[Event, EventSourcing]) SaveSnapshotAndEvent2(ctx context.Context, aggregates []EventSourcing) error {
	es := aggregates[len(aggregates)-1].DomainEvents()
	events := make([]Event, 0)
	for i := 0; i < len(es); i++ {
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

	s.eventPulisher.Publish(ctx, es)

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

func (s *ApplicationService[Event, EventSourcing]) withUnitOfWork(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.innerService.withUnitOfWork(ctx, fn)
}
