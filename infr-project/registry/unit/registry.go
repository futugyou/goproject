package unit

import (
	"context"
	"fmt"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/options"
)

type Registry struct {
	Options *options.Options
	events  map[string]func(context.Context, options.Options) domain.IUnitOfWork
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		Options: &options.Options{},
		events:  map[string]func(context.Context, options.Options) domain.IUnitOfWork{},
	}
}

func (s *Registry) RegisterComponent(
	componentFactory func(context.Context, options.Options) domain.IUnitOfWork,
	names ...string,
) {
	for _, name := range names {
		s.events[fmt.Sprintf("unit-of-work-%s", name)] = componentFactory
	}
}

func (s *Registry) Create(ctx context.Context) (domain.IUnitOfWork, error) {
	if s.Options == nil {
		return nil, fmt.Errorf("options is nil")
	}

	name := s.Options.StoreType
	if method, ok := s.events[fmt.Sprintf("unit-of-work-%s", name)]; ok {
		return method(ctx, *s.Options), nil
	}

	return nil, fmt.Errorf("couldn't find unit of work %s", name)
}
