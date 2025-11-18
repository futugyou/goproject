package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/platformservice/infrastructure/internal/entity"
)

type EventLogDao struct {
	mongoimpl.BaseCRUD[entity.EventLogEntity]
}

func NewEventLogDao(client *mongo.Client, config mongoimpl.DBConfig) *EventLogDao {
	if config.CollectionName == "" {
		config.CollectionName = "platform_event_logs"
	}

	getID := func(a entity.EventLogEntity) string { return a.ID }

	return &EventLogDao{
		BaseCRUD: *mongoimpl.NewBaseCRUD(client, config, getID),
	}
}

func (s *EventLogDao) InsertMany(ctx context.Context, entities []entity.EventLogEntity) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)
	documents := make([]any, 0, len(entities))
	for _, entity := range entities {
		documents = append(documents, entity)
	}

	_, err := c.InsertMany(ctx, documents)
	return err
}
