package component

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/options"
	"github.com/futugyou/infr-project/registry/unit"
)

func init() {
	unit.DefaultRegistry.RegisterComponent(func(ctx context.Context, option options.Options) domain.IUnitOfWork {
		config := infrastructure_mongo.DBConfig{
			ConnectString: option.MongoDBURL,
		}
		mongoclient, _ := mongo.Connect(ctx, mongo_options.Client().ApplyURI(config.ConnectString))
		u, _ := infrastructure_mongo.NewMongoUnitOfWork(mongoclient)
		return u
	}, func(ctx context.Context, option options.Options) domain.IUnitOfWork {
		config := infrastructure_mongo.DBConfig{
			ConnectString: option.MongoDBURL,
		}
		mongoclient, _ := mongo.Connect(ctx, mongo_options.Client().ApplyURI(config.ConnectString))
		u, _ := infrastructure_mongo.NewMongoUnitOfWork(mongoclient)
		return u
	}, "mongo")
}
