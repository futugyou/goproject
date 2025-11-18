package dao

import (
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
