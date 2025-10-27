package unit

import (
	"context"
	"fmt"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/options"
)

type Registry struct {
	Options     *options.Options
	events      map[string]func(context.Context, options.Options) domain.IUnitOfWork
	eventAsyncs map[string]func(context.Context, options.Options) domain.IUnitOfWork
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		Options:     &options.Options{},
		events:      map[string]func(context.Context, options.Options) domain.IUnitOfWork{},
		eventAsyncs: map[string]func(context.Context, options.Options) domain.IUnitOfWork{},
	}
}

func (s *Registry) RegisterComponent(
	componentFactory func(context.Context, options.Options) domain.IUnitOfWork,
	componentAsyncFactory func(context.Context, options.Options) domain.IUnitOfWork,
	names ...string,
) {
	for _, name := range names {
		s.events[fmt.Sprintf("unit-of-work-%s", name)] = componentFactory
		s.eventAsyncs[fmt.Sprintf("unit-of-work-async-%s", name)] = componentAsyncFactory
	}
}

func (s *Registry) Create(ctx context.Context) (domain.IUnitOfWork, domain.IUnitOfWork, error) {
	if s.Options == nil {
		return nil, nil, fmt.Errorf("options is nil")
	}

	name := s.Options.StoreType
	if method, ok := s.events[fmt.Sprintf("unit-of-work-%s", name)]; ok {
		if methodasync, ok2 := s.eventAsyncs[fmt.Sprintf("unit-of-work-async-%s", name)]; ok2 {
			return method(ctx, *s.Options), methodasync(ctx, *s.Options), nil
		}
	}

	return nil, nil, fmt.Errorf("couldn't find unit of work %s", name)
}
