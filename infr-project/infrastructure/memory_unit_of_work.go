package infrastructure

import (
	"github.com/futugyou/infr-project/domain"
)

type MemoryUnitOfWork struct {
	repository domain.IEventSourcingRepository[domain.IEventSourcing]
	events     []domain.IDomainEvent
	aggregate  domain.IEventSourcing
}

func NewMemoryUnitOfWork(
	repository domain.IEventSourcingRepository[domain.IEventSourcing],
) *MemoryUnitOfWork {
	return &MemoryUnitOfWork{
		repository: repository,
	}
}

func (u *MemoryUnitOfWork) RegisterNew(aggregate domain.IEventSourcing, events []domain.IDomainEvent) {
	u.aggregate = aggregate
	u.events = append(u.events, events...)
}

func (u *MemoryUnitOfWork) Commit() error {
	// Save the aggregate and its events
	if err := u.repository.Save(u.aggregate); err != nil {
		return err
	}

	// Clear the events and aggregate references after commit
	u.events = nil
	u.aggregate = nil

	return nil
}

func (u *MemoryUnitOfWork) Rollback() error {
	// In-memory rollback might not need any action
	u.events = nil
	u.aggregate = nil
	return nil
}
