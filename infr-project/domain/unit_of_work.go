package domain

import (
	"context"
	"fmt"

	"github.com/futugyou/infr-project/options"
)

type IUnitOfWork interface {
	Start(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type IUnitOfWorkAsync interface {
	StartAsync(ctx context.Context) (<-chan context.Context, <-chan error)
	CommitAsync(ctx context.Context) <-chan error
	RollbackAsync(ctx context.Context) <-chan error
}

type UnitOfWorkRegistry struct {
	Options     *options.Options
	events      map[string]func(options.Options) IUnitOfWork
	eventAsyncs map[string]func(options.Options) IUnitOfWorkAsync
}

var DefaultUnitOfWorkRegistry *UnitOfWorkRegistry = NewUnitOfWorkRegistry()

func NewUnitOfWorkRegistry() *UnitOfWorkRegistry {
	return &UnitOfWorkRegistry{
		Options:     &options.Options{},
		events:      map[string]func(options.Options) IUnitOfWork{},
		eventAsyncs: map[string]func(options.Options) IUnitOfWorkAsync{},
	}
}

func (s *UnitOfWorkRegistry) RegisterComponent(componentFactory func(options.Options) IUnitOfWork, componentAsyncFactory func(options.Options) IUnitOfWorkAsync, names ...string) {
	for _, name := range names {
		s.events[fmt.Sprintf("unit-of-work-%s", name)] = componentFactory
		s.eventAsyncs[fmt.Sprintf("unit-of-work-async-%s", name)] = componentAsyncFactory
	}
}

func (s *UnitOfWorkRegistry) Create() (IUnitOfWork, IUnitOfWorkAsync, error) {
	if s.Options == nil {
		return nil, nil, fmt.Errorf("options is nil")
	}

	name := s.Options.StoreType
	if method, ok := s.events[fmt.Sprintf("unit-of-work-%s", name)]; ok {
		if methodasync, ok2 := s.eventAsyncs[fmt.Sprintf("unit-of-work-async-%s", name)]; ok2 {
			return method(*s.Options), methodasync(*s.Options), nil
		}
	}

	return nil, nil, fmt.Errorf("couldn't find unit of work %s", name)
}
