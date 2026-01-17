package component

import (
	"context"

	"github.com/futugyou/domaincore/application"
	"github.com/futugyou/domaincore/qstashdispatcherimpl"
	"github.com/futugyou/infr-project/registry/options"
	"github.com/futugyou/infr-project/registry/publisher"
)

func init() {
	publisher.DefaultRegistry.RegisterComponent(func(ctx context.Context, option options.Options) application.EventDispatcher {
		return qstashdispatcherimpl.NewQStashEventDispatcher(option.QstashToken, option.QstashDestination)
	}, "qstash")
}
