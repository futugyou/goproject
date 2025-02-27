package component

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/options"
	"github.com/futugyou/infr-project/registry/snapshot_store"
)

func init() {
	snapshot_store.DefaultRegistry.RegisterComponent(func(ctx context.Context, option options.Options) infrastructure.ISnapshotStore[domain.IEventSourcing] {
		config := infrastructure_mongo.DBConfig{
			DBName:        option.DBName,
			ConnectString: option.MongoDBURL,
		}
		client, _ := mongo.Connect(ctx, mongo_options.Client().ApplyURI(config.ConnectString))
		return infrastructure_mongo.NewMongoSnapshotStore[domain.IEventSourcing](client, config)
	}, func(ctx context.Context, option options.Options) infrastructure.ISnapshotStoreAsync[domain.IEventSourcing] {
		config := infrastructure_mongo.DBConfig{
			DBName:        option.DBName,
			ConnectString: option.MongoDBURL,
		}
		client, _ := mongo.Connect(ctx, mongo_options.Client().ApplyURI(config.ConnectString))
		return infrastructure_mongo.NewMongoSnapshotStore[domain.IEventSourcing](client, config)
	}, "mongo")
}
