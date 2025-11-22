package unit

import (
	"context"
	"fmt"

	"github.com/futugyou/domaincore/domain"
	"github.com/futugyou/infr-project/registry/options"
)

type Registry struct {
	Options *options.Options
	events  map[string]func(context.Context, options.Options) domain.UnitOfWork
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		Options: &options.Options{},
		events:  map[string]func(context.Context, options.Options) domain.UnitOfWork{},
	}
}

func (s *Registry) RegisterComponent(
	componentFactory func(context.Context, options.Options) domain.UnitOfWork,
	names ...string,
) {
	for _, name := range names {
		s.events[fmt.Sprintf("unit-of-work-%s", name)] = componentFactory
	}
}

func (s *Registry) Create(ctx context.Context) (domain.UnitOfWork, error) {
	if s.Options == nil {
		return nil, fmt.Errorf("options is nil")
	}

	name := s.Options.StoreType
	if method, ok := s.events[fmt.Sprintf("unit-of-work-%s", name)]; ok {
		return method(ctx, *s.Options), nil
	}

	return nil, fmt.Errorf("couldn't find unit of work %s", name)
}
