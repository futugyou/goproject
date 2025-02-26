package publisher

import (
	"fmt"

	"github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/options"
)

type Registry struct {
	Options    *options.Options
	publishers map[string]func(options.Options) infrastructure.IEventPublisher
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		publishers: map[string]func(options.Options) infrastructure.IEventPublisher{},
	}
}

func (s *Registry) RegisterComponent(componentFactory func(options.Options) infrastructure.IEventPublisher, names ...string) {
	for _, name := range names {
		s.publishers[fmt.Sprintf("event-publisher-%s", name)] = componentFactory
	}
}

func (s *Registry) Create() (infrastructure.IEventPublisher, error) {
	if s.Options == nil {
		return nil, fmt.Errorf("options is nil")
	}

	name := s.Options.EventPublisher
	if method, ok := s.publishers[fmt.Sprintf("event-publisher-%s", name)]; ok {
		return method(*s.Options), nil
	}
	return nil, fmt.Errorf("couldn't find event publisher %s", name)
}
