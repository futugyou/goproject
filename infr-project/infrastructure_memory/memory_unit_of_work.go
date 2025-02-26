package infrastructure_memory

import (
	"context"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/options"
)

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

func (u *MemoryUnitOfWork) StartAsync(ctx context.Context) (<-chan context.Context, <-chan error) {
	ctxChan := make(chan context.Context, 1)
	errChan := make(chan error, 1)
	ctxChan <- ctx
	close(ctxChan)
	errChan <- nil
	close(errChan)
	return ctxChan, errChan
}

func (u *MemoryUnitOfWork) CommitAsync(ctx context.Context) <-chan error {
	errChan := make(chan error, 1)
	errChan <- nil
	close(errChan)
	return errChan
}

func (u *MemoryUnitOfWork) RollbackAsync(ctx context.Context) <-chan error {
	errChan := make(chan error, 1)
	errChan <- nil
	close(errChan)
	return errChan
}

func init() {
	domain.DefaultUnitOfWorkRegistry.RegisterComponent(func(option options.Options) domain.IUnitOfWork {
		return NewMemoryUnitOfWork()
	}, func(option options.Options) domain.IUnitOfWorkAsync {
		return NewMemoryUnitOfWork()
	}, "memory")
}
