package component

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/options"
	"github.com/futugyou/infr-project/registry/event_store"
	"github.com/futugyou/infr-project/resource"
)

func init() {
	create := func(eventType string) (domain.IDomainEvent, error) {
		return resource.CreateEvent(eventType)
	}
	event_store.DefaultRegistry.RegisterComponent(func(option options.Options) infrastructure.IEventStore[domain.IDomainEvent] {
		ctx := context.Background()
		config := infrastructure_mongo.DBConfig{
			DBName:        option.DBName,
			ConnectString: option.MongoDBURL,
		}
		client, _ := mongo.Connect(ctx, mongo_options.Client().ApplyURI(config.ConnectString))
		return infrastructure_mongo.NewMongoEventStore(client, config, "resource_events", create)
	}, func(option options.Options) infrastructure.IEventStoreAsync[domain.IDomainEvent] {
		ctx := context.Background()
		config := infrastructure_mongo.DBConfig{
			DBName:        option.DBName,
			ConnectString: option.MongoDBURL,
		}
		client, _ := mongo.Connect(ctx, mongo_options.Client().ApplyURI(config.ConnectString))
		return infrastructure_mongo.NewMongoEventStore(client, config, "resource_events", create)
	}, "memory")
}
