package component

import (
	"context"

	"github.com/futugyou/domaincore/infrastructure"
	"github.com/futugyou/domaincore/qstashdispatcherimpl"
	"github.com/futugyou/infr-project/registry/options"
	"github.com/futugyou/infr-project/registry/publisher"
)

func init() {
	publisher.DefaultRegistry.RegisterComponent(func(ctx context.Context, option options.Options) infrastructure.EventDispatcher {
		return qstashdispatcherimpl.NewQStashEventDispatcher(option.QstashToken, option.QstashDestination)
	}, "qstash")
}
