package component

import (
	"context"

	"github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/infrastructure_qstash"
	"github.com/futugyou/infr-project/options"
	"github.com/futugyou/infr-project/registry/publisher"
)

func init() {
	publisher.DefaultRegistry.RegisterComponent(func(ctx context.Context, option options.Options) infrastructure.IEventPublisher {
		return infrastructure_qstash.NewQStashEventPulisher(option.QstashToken, option.QstashDestination)
	}, "qstash")
}
