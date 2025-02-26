package screenshot

import (
	"fmt"

	"github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/options"
)

type Registry struct {
	Options *options.Options
	events  map[string]func(options.Options) infrastructure.IScreenshot
}

var DefaultRegistry *Registry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		events: map[string]func(options.Options) infrastructure.IScreenshot{},
	}
}

func (s *Registry) RegisterComponent(componentFactory func(options.Options) infrastructure.IScreenshot, names ...string) {
	for _, name := range names {
		s.events[fmt.Sprintf("screenshot-%s", name)] = componentFactory
	}
}

func (s *Registry) Create() (infrastructure.IScreenshot, error) {
	if s.Options == nil {
		return nil, fmt.Errorf("options is nil")
	}

	name := s.Options.ScreenshotType
	if method, ok := s.events[fmt.Sprintf("screenshot-%s", name)]; ok {
		return method(*s.Options), nil
	}
	return nil, fmt.Errorf("couldn't find screenshot %s", name)
}
