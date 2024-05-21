package infrastructure_memory

import "context"

type MemoryUnitOfWork struct {
}

func NewMemoryUnitOfWork() *MemoryUnitOfWork {
	return &MemoryUnitOfWork{}
}

func (u *MemoryUnitOfWork) Start(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (u *MemoryUnitOfWork) Commit(ctx context.Context) error {
	return nil
}

func (u *MemoryUnitOfWork) Rollback(ctx context.Context) error {
	return nil
}

func (u *MemoryUnitOfWork) End(ctx context.Context) {}
