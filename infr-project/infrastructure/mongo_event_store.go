package infrastructure

import (
	"context"

	"github.com/futugyou/infr-project/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConfig struct {
	DBName        string
	ConnectString string
}

type MongoEventStore[Event domain.IDomainEvent] struct {
	DBName string
	Client *mongo.Client
}

func NewMongoEventStore[Event domain.IDomainEvent](config DBConfig) *MongoEventStore[Event] {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	return &MongoEventStore[Event]{
		DBName: config.DBName,
		Client: client,
	}
}

func (s *MongoEventStore[Event]) Load(id string) ([]Event, error) {
	c := s.Client.Database(s.DBName).Collection("domain_events")
	result := make([]Event, 0)

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

func (s *MongoEventStore[Event]) Save(events []Event) error {
	if len(events) == 0 {
		return nil
	}

	c := s.Client.Database(s.DBName).Collection("domain_events")
	entitys := make([]interface{}, len(events))
	for i := 0; i < len(events); i++ {
		entitys[i] = events[i]
	}

	ctx := context.Background()
	_, err := c.InsertMany(ctx, entitys)

	return err
}
