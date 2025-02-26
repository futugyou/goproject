package unit

import (
	"fmt"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/options"
)

type Registry struct {
	Options     *options.Options
	events      map[string]func(options.Options) domain.IUnitOfWork
	eventAsyncs map[string]func(options.Options) domain.IUnitOfWorkAsync
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		Options:     &options.Options{},
		events:      map[string]func(options.Options) domain.IUnitOfWork{},
		eventAsyncs: map[string]func(options.Options) domain.IUnitOfWorkAsync{},
	}
}

func (s *Registry) RegisterComponent(componentFactory func(options.Options) domain.IUnitOfWork, componentAsyncFactory func(options.Options) domain.IUnitOfWorkAsync, names ...string) {
	for _, name := range names {
		s.events[fmt.Sprintf("unit-of-work-%s", name)] = componentFactory
		s.eventAsyncs[fmt.Sprintf("unit-of-work-async-%s", name)] = componentAsyncFactory
	}
}

func (s *Registry) Create() (domain.IUnitOfWork, domain.IUnitOfWorkAsync, error) {
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
