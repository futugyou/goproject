package infrastructure_mongo

import (
	"context"
	"fmt"

	"github.com/chidiwilliams/flatbson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/domain"
)

type DBConfig struct {
	DBName        string
	ConnectString string
}

type MongoEventStore[Event domain.IDomainEvent] struct {
	DBName       string
	Client       *mongo.Client
	TableName    string
	EventFactory func(eventType string) (Event, error)
}

func NewMongoEventStore[Event domain.IDomainEvent](
	client *mongo.Client,
	config DBConfig,
	tableName string,
	eventFactory func(eventType string) (Event, error),
) *MongoEventStore[Event] {
	return &MongoEventStore[Event]{
		DBName:       config.DBName,
		Client:       client,
		TableName:    tableName,
		EventFactory: eventFactory,
	}
}

func (s *MongoEventStore[Event]) Load(ctx context.Context, id string) ([]Event, error) {
	filter := bson.D{{Key: "id", Value: id}}
	return s.load(ctx, filter)
}

func (s *MongoEventStore[Event]) LoadGreaterthanVersion(ctx context.Context, id string, version int) ([]Event, error) {
	filter := bson.D{
		{Key: "id", Value: id},
		{Key: "version", Value: bson.D{
			{Key: "$gt", Value: version},
		}},
	}
	return s.load(ctx, filter)
}

func (s *MongoEventStore[Event]) load(ctx context.Context, filter primitive.D) ([]Event, error) {
	c := s.Client.Database(s.DBName).Collection(s.TableName)
	result := make([]Event, 0)
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var rawDocs []bson.Raw
	if err = cursor.All(ctx, &rawDocs); err != nil {
		return nil, err
	}

	for _, rawDoc := range rawDocs {
		var eventType struct {
			Type string `bson:"event_type"`
		}
		if err := bson.Unmarshal(rawDoc, &eventType); err != nil {
			return nil, err
		}

		event, err := s.EventFactory(eventType.Type)
		if err != nil {
			return nil, err
		}

		if err := bson.Unmarshal(rawDoc, event); err != nil {
			return nil, err
		}

		result = append(result, event)
	}

	return result, nil
}

func (s *MongoEventStore[Event]) Save(ctx context.Context, events []Event) (err error) {
	if len(events) == 0 {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	c := s.Client.Database(s.DBName).Collection(s.TableName)
	entities := make([]interface{}, len(events))
	for i := 0; i < len(events); i++ {
		eventType := events[i].EventType()
		eventMap, err := flatbson.Flatten(events[i])
		if err != nil {
			return err
		}

		eventMap["event_type"] = eventType
		eventMap["is_handled"] = 0
		entities[i] = eventMap
	}

	_, err = c.InsertMany(ctx, entities)

	return err
}
