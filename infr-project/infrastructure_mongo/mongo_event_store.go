package infrastructure_mongo

import (
	"context"

	"github.com/chidiwilliams/flatbson"
	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/resource"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBConfig struct {
	DBName        string
	ConnectString string
}

type MongoEventStore[Event domain.IDomainEvent] struct {
	DBName string
	Client *mongo.Client
}

func NewMongoEventStore[Event domain.IDomainEvent](client *mongo.Client, config DBConfig) *MongoEventStore[Event] {
	return &MongoEventStore[Event]{
		DBName: config.DBName,
		Client: client,
	}
}

func (s *MongoEventStore[Event]) Load(id string) ([]Event, error) {
	filter := bson.D{{Key: "id", Value: id}}
	return s.load(filter)
}

func (s *MongoEventStore[Event]) LoadGreaterthanVersion(id string, version int) ([]Event, error) {
	filter := bson.D{
		{Key: "id", Value: id},
		{Key: "version", Value: bson.D{
			{Key: "$gt", Value: version},
		}},
	}
	return s.load(filter)
}

func (s *MongoEventStore[Event]) load(filter primitive.D) ([]Event, error) {
	e := new(Event)
	c := s.Client.Database(s.DBName).Collection((*e).AggregateEventName())
	result := make([]Event, 0)
	ctx := context.Background()
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

		event, err := resource.CreateEvent(eventType.Type)
		if err != nil {
			return nil, err
		}

		if err := bson.Unmarshal(rawDoc, event); err != nil {
			return nil, err
		}

		result = append(result, event.(Event))
	}

	return result, nil
}

func (s *MongoEventStore[Event]) Save(ctx context.Context, events []Event) error {
	if len(events) == 0 {
		return nil
	}

	e := new(Event)
	c := s.Client.Database(s.DBName).Collection((*e).AggregateEventName())
	entities := make([]interface{}, len(events))
	for i := 0; i < len(events); i++ {
		eventType := events[i].EventType()
		eventMap, err := flatbson.Flatten(events[i])
		if err != nil {
			return err
		}

		eventMap["event_type"] = eventType
		entities[i] = eventMap
	}

	_, err := c.InsertMany(ctx, entities)

	return err
}
