package infrastructure

import (
	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/resourceservice/domain"
	"github.com/futugyou/resourceservice/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResourceSnapshotStore struct {
	*mongoimpl.MongoSnapshotStore[*domain.Resource]
}

func NewResourceSnapshotStore(mongoclient *mongo.Client, option *options.Options) *ResourceSnapshotStore {
	snapshotStoreConfig := mongoimpl.DBConfig{
		DBName:         option.DBName,
		CollectionName: "resources",
	}
	snapshotStore := mongoimpl.NewMongoSnapshotStore[*domain.Resource](mongoclient, snapshotStoreConfig)
	return &ResourceSnapshotStore{
		MongoSnapshotStore: snapshotStore,
	}
}
