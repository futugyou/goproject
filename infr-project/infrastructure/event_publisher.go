package infrastructure

import (
	"context"
	"fmt"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/options"
)

type IEventPublisher interface {
	Publish(ctx context.Context, events []domain.IDomainEvent) error
	PublishCommon(ctx context.Context, event any, event_type string) error
}

type EventPublisherRegistry struct {
	Options    *options.Options
	publishers map[string]func(options.Options) IEventPublisher
}

var DefaultEventPublisherRegistry *EventPublisherRegistry = NewEventPublisherRegistry()

func NewEventPublisherRegistry() *EventPublisherRegistry {
	return &EventPublisherRegistry{
		publishers: map[string]func(options.Options) IEventPublisher{},
	}
}

func (s *EventPublisherRegistry) RegisterComponent(componentFactory func(options.Options) IEventPublisher, names ...string) {
	for _, name := range names {
		s.publishers[fmt.Sprintf("event-publisher-%s", name)] = componentFactory
	}
}

func (s *EventPublisherRegistry) Create() (IEventPublisher, error) {
	if s.Options == nil {
		return nil, fmt.Errorf("options is nil")
	}

	name := s.Options.EventPublisher
	if method, ok := s.publishers[fmt.Sprintf("event-publisher-%s", name)]; ok {
		return method(*s.Options), nil
	}
	return nil, fmt.Errorf("couldn't find event publisher %s", name)
}
