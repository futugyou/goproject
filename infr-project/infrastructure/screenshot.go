package infrastructure

import (
	"context"
	"fmt"

	"github.com/futugyou/infr-project/options"
)

type IScreenshot interface {
	Create(ctx context.Context, url string) (*string, error)
}

type ScreenshotRegistry struct {
	Options *options.Options
	events  map[string]func(options.Options) IScreenshot
}

var DefaultScreenshotRegistry *ScreenshotRegistry = NewScreenshotRegistry()

func NewScreenshotRegistry() *ScreenshotRegistry {
	return &ScreenshotRegistry{
		events: map[string]func(options.Options) IScreenshot{},
	}
}

func (s *ScreenshotRegistry) RegisterComponent(componentFactory func(options.Options) IScreenshot, names ...string) {
	for _, name := range names {
		s.events[fmt.Sprintf("screenshot-%s", name)] = componentFactory
	}
}

func (s *ScreenshotRegistry) Create() (IScreenshot, error) {
	if s.Options == nil {
		return nil, fmt.Errorf("options is nil")
	}

	name := s.Options.ScreenshotType
	if method, ok := s.events[fmt.Sprintf("screenshot-%s", name)]; ok {
		return method(*s.Options), nil
	}
	return nil, fmt.Errorf("couldn't find screenshot %s", name)
}
