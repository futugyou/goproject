package infrastructure_memory

import (
	"context"
	"fmt"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/options"
)

type MemoryEventStore[Event domain.IDomainEvent] struct {
	storage map[string][]Event
}

func NewMemoryEventStore[Event domain.IDomainEvent]() *MemoryEventStore[Event] {
	return &MemoryEventStore[Event]{
		storage: make(map[string][]Event),
	}
}

func (s *MemoryEventStore[Event]) LoadGreaterthanVersion(ctx context.Context, id string, version int) ([]Event, error) {
	events, ok := s.storage[id]
	if !ok {
		return nil, fmt.Errorf("no data for %s", id)
	}

	result := make([]Event, 0)
	for i := 0; i < len(events); i++ {
		if events[i].Version() > version {
			result = append(result, events[i])
		}
	}

	return result, nil
}

func (s *MemoryEventStore[Event]) Load(ctx context.Context, id string) ([]Event, error) {
	events, ok := s.storage[id]
	if !ok {
		return nil, fmt.Errorf("no data for %s", id)
	}

	return events, nil
}

func (s *MemoryEventStore[Event]) Save(ctx context.Context, events []Event) error {
	for _, event := range events {
		id := event.AggregateId()
		s.storage[id] = append(s.storage[id], event)
	}
	return nil
}

func (s *MemoryEventStore[Event]) LoadGreaterthanVersionAsync(ctx context.Context, id string, version int) (<-chan []Event, <-chan error) {
	resultChan := make(chan []Event, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		result, err := s.LoadGreaterthanVersion(ctx, id, version)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	return resultChan, errorChan
}

func (s *MemoryEventStore[Event]) LoadAsync(ctx context.Context, id string) (<-chan []Event, <-chan error) {
	resultChan := make(chan []Event, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		result, err := s.Load(ctx, id)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	return resultChan, errorChan
}

func (s *MemoryEventStore[Event]) SaveAsync(ctx context.Context, events []Event) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)

		err := s.Save(ctx, events)
		errorChan <- err
	}()

	return errorChan
}

func init() {
	infrastructure.DefaultEventStoreRegistry.RegisterComponent(func(option options.Options) infrastructure.IEventStore[domain.IDomainEvent] {
		return NewMemoryEventStore[domain.IDomainEvent]()
	}, func(option options.Options) infrastructure.IEventStoreAsync[domain.IDomainEvent] {
		return NewMemoryEventStore[domain.IDomainEvent]()
	}, "memory")
}
