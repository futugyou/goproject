package infrastructure

import (
	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/resourceservice/domain"
	"github.com/futugyou/resourceservice/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResourceEventStore struct {
	*mongoimpl.MongoEventStore[domain.ResourceEvent]
}

func NewResourceEventStore(mongoclient *mongo.Client, option *options.Options) *ResourceEventStore {
	eventStoreConfig := mongoimpl.DBConfig{
		DBName:         option.DBName,
		CollectionName: "resource_events",
	}
	eventStore := mongoimpl.NewMongoEventStore(mongoclient, eventStoreConfig, domain.CreateEvent)

	return &ResourceEventStore{
		MongoEventStore: eventStore,
	}
}
