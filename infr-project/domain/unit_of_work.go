package domain

import "context"

type IUnitOfWork interface {
	Start(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	End(ctx context.Context)
}

type IUnitOfWorkAsync interface {
	StartAsync(ctx context.Context) (<-chan context.Context, <-chan error)
	CommitAsync(ctx context.Context) <-chan error
	RollbackAsync(ctx context.Context) <-chan error
	EndAsync(ctx context.Context) <-chan error
}
