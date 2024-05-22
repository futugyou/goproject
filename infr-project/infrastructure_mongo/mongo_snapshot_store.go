package infrastructure_mongo

import (
	"context"

	"github.com/futugyou/infr-project/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSnapshotStore[EventSourcing domain.IEventSourcing] struct {
	DBName string
	Client *mongo.Client
}

func NewMongoSnapshotStore[EventSourcing domain.IEventSourcing](client *mongo.Client, config DBConfig) *MongoSnapshotStore[EventSourcing] {
	return &MongoSnapshotStore[EventSourcing]{
		DBName: config.DBName,
		Client: client,
	}
}

func (s *MongoSnapshotStore[EventSourcing]) LoadSnapshot(id string) ([]EventSourcing, error) {
	a := new(EventSourcing)
	c := s.Client.Database(s.DBName).Collection((*a).AggregateName())
	result := make([]EventSourcing, 0)

	filter := bson.D{{Key: "id", Value: id}}
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

// func (s *MongoSnapshotStore[EventSourcing]) LoadLatestSnapshot(id string) (*EventSourcing, error) {
// 	a := new(EventSourcing)
// 	c := s.Client.Database(s.DBName).Collection((*a).AggregateName())

// 	filter := bson.D{{Key: "id", Value: id}}
// 	opts := &options.FindOneOptions{}
// 	opts.SetSort(bson.D{{Key: "version", Value: -1}})
// 	ctx := context.Background()
// 	if err := c.FindOne(ctx, filter, opts).Decode(&a); err != nil {
// 		return nil, err
// 	}

// 	return a, nil
// }

func (s *MongoSnapshotStore[EventSourcing]) SaveSnapshot(ctx context.Context, aggregate EventSourcing) error {
	c := s.Client.Database(s.DBName).Collection(aggregate.AggregateName())
	_, err := c.InsertOne(ctx, aggregate)
	return err
}
