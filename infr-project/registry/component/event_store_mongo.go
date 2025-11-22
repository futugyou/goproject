package component

import (
	"context"

	coredomain "github.com/futugyou/domaincore/domain"
	coreinfr "github.com/futugyou/domaincore/infrastructure"
	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/infr-project/registry/event_store"
	"github.com/futugyou/infr-project/registry/options"
)

func init() {
	event_store.DefaultRegistry.RegisterComponent(func(ctx context.Context, option options.Options) coreinfr.EventStore[coredomain.DomainEvent] {
		mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.MongoDBURL)
		if err != nil {
			return nil
		}

		eventStoreConfig := mongoimpl.DBConfig{
			DBName:         option.DBName,
			CollectionName: "resource_events",
		}

		return mongoimpl.NewMongoEventStore(mongoclient, eventStoreConfig, option.EventFactory)
	}, "mongo")
}
