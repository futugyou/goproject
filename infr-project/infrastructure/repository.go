package infrastructure

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

// Deprecated: EventSourcingRepository is deprecated.
type EventSourcingRepository[Aggregate domain.IEventSourcing] struct {
	eventStore        IEventStore[domain.IDomainEvent]
	snapshotStore     ISnapshotStore[Aggregate]
	needStoreSnapshot func(Aggregate) bool
	newAggregateFunc  func() Aggregate
}

// Deprecated: EventSourcingRepository is deprecated.
func NewEventSourcingRepository[Aggregate domain.IEventSourcing](
	eventStore IEventStore[domain.IDomainEvent],
	snapshotStore ISnapshotStore[Aggregate],
	needStoreSnapshot func(Aggregate) bool,
	newAggregateFunc func() Aggregate,
) *EventSourcingRepository[Aggregate] {
	return &EventSourcingRepository[Aggregate]{
		eventStore:        eventStore,
		snapshotStore:     snapshotStore,
		needStoreSnapshot: needStoreSnapshot,
		newAggregateFunc:  newAggregateFunc,
	}
}

func (r *EventSourcingRepository[Aggregate]) Load(id string) (*Aggregate, error) {
	// Attempt to restore from the latest snapshot
	datas, err := r.snapshotStore.LoadSnapshot(id)
	// we save snapshot from fisrt version(1), so it contains at least one piece of data
	if err != nil || len(datas) == 0 {
		return nil, err
	}

	aggregate := datas[len(datas)-1]

	// Load all events for the aggregate
	events, err := r.eventStore.Load(id)
	if err != nil {
		return nil, err
	}

	// Apply events to the aggregate to restore its state
	for i := 0; i < len(events); i++ {
		if aggregate.AggregateVersion() >= events[i].Version() {
			continue
		}
		aggregate.Apply(events[i])
	}

	return &aggregate, nil
}

func (r *EventSourcingRepository[Aggregate]) LoadAll(id string) ([]Aggregate, error) {
	events, err := r.eventStore.Load(id)
	if err != nil {
		return nil, err
	}

	aggregate := r.newAggregateFunc()
	result := make([]Aggregate, 0)
	// Apply events to the aggregate to restore its state
	for i := 0; i < len(events); i++ {
		aggregate.Apply(events[i])
		result = append(result, aggregate)
	}

	return result, nil
}

func (r *EventSourcingRepository[Aggregate]) Save(aggregate Aggregate) error {
	// Save the events
	if err := r.eventStore.Save(context.Background(), aggregate.DomainEvents()); err != nil {
		return err
	}

	// Take a snapshot if necessary
	if r.needStoreSnapshot(aggregate) {
		if err := r.snapshotStore.SaveSnapshot(context.Background(), aggregate); err != nil {
			return err
		}
	}

	// Clear the uncommitted events
	aggregate.ClearDomainEvents()
	return nil
}
