package mongoimpl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/domaincore/domain"
)

type MongoSnapshotStore[EventSourcing domain.EventSourcing] struct {
	DBName         string
	CollectionName string
	Client         *mongo.Client
}

func NewMongoSnapshotStore[EventSourcing domain.EventSourcing](client *mongo.Client, config DBConfig) *MongoSnapshotStore[EventSourcing] {
	collectionName := config.CollectionName
	if collectionName == "" {
		collectionName = (*new(EventSourcing)).AggregateName()
	}
	return &MongoSnapshotStore[EventSourcing]{
		CollectionName: collectionName,
		DBName:         config.DBName,
		Client:         client,
	}
}

func (s *MongoSnapshotStore[EventSourcing]) LoadSnapshot(ctx context.Context, id string) ([]EventSourcing, error) {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)
	result := make([]EventSourcing, 0)

	filter := bson.D{{Key: "id", Value: id}}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result, nil
}

func (s *MongoSnapshotStore[EventSourcing]) SaveSnapshot(ctx context.Context, aggregate EventSourcing) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)
	_, err := c.InsertOne(ctx, aggregate)
	return err
}
