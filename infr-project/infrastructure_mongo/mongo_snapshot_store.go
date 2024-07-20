package infrastructure_mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/domain"
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

func (s *MongoSnapshotStore[EventSourcing]) LoadSnapshot(ctx context.Context, id string) ([]EventSourcing, error) {
	a := new(EventSourcing)
	c := s.Client.Database(s.DBName).Collection((*a).AggregateName())
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
	c := s.Client.Database(s.DBName).Collection(aggregate.AggregateName())
	_, err := c.InsertOne(ctx, aggregate)
	return err
}

func (s *MongoSnapshotStore[EventSourcing]) LoadSnapshotAsync(ctx context.Context, id string) (<-chan []EventSourcing, <-chan error) {
	resultChan := make(chan []EventSourcing, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		result, err := s.LoadSnapshot(ctx, id)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	return resultChan, errorChan
}

func (s *MongoSnapshotStore[EventSourcing]) SaveSnapshotAsync(ctx context.Context, aggregate EventSourcing) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)

		err := s.SaveSnapshot(ctx, aggregate)
		errorChan <- err
	}()

	return errorChan
}
