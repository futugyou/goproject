package infrastructure

import (
	"context"

	"github.com/futugyou/infr-project/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSnapshotStore[EventSourcing domain.IEventSourcing] struct {
	DBName string
	Client *mongo.Client
}

func NewMongoSnapshotStore[EventSourcing domain.IEventSourcing](config DBConfig) *MongoSnapshotStore[EventSourcing] {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	return &MongoSnapshotStore[EventSourcing]{
		DBName: config.DBName,
		Client: client,
	}
}

func (s *MongoSnapshotStore[EventSourcing]) LoadSnapshot(id string) ([]EventSourcing, error) {
	a := new(EventSourcing)
	c := s.Client.Database(s.DBName).Collection((*a).AggregateName())
	result := make([]EventSourcing, 0)

	filter := bson.D{{Key: "Id", Value: id}}
	ctx := context.Background()
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

func (s *MongoSnapshotStore[EventSourcing]) SaveSnapshot(aggregate EventSourcing) error {
	c := s.Client.Database(s.DBName).Collection(aggregate.AggregateName())
	ctx := context.Background()
	_, err := c.InsertOne(ctx, aggregate)
	return err
}
