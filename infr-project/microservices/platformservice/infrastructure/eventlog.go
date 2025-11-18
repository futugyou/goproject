package infrastructure

import (
	"context"

	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/platformservice/application/service"
	"github.com/futugyou/platformservice/infrastructure/internal/dao"
	"github.com/futugyou/platformservice/infrastructure/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventLogService struct {
	log *dao.EventLogDao
}

func NewEventLogService(client *mongo.Client, config mongoimpl.DBConfig) *EventLogService {
	log := dao.NewEventLogDao(client, config)

	return &EventLogService{
		log: log,
	}
}

// Create implements service.EventLogService.
func (e *EventLogService) Create(ctx context.Context, event service.EventLog) error {
	entity := entity.EventLogEntity{
		ID:         event.EventID,
		PlatformID: event.PlatformID,
		ProjectID:  event.ProjectID,
		Token:      event.Token,
		EventType:  event.EventType,
		CreatedAt:  event.CreatedAt,
	}

	return e.log.Insert(ctx, entity)
}

// Get implements service.EventLogService.
func (e *EventLogService) Get(ctx context.Context, eventID string) (*service.EventLog, error) {
	entity, err := e.log.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return &service.EventLog{
		EventID:    entity.ID,
		PlatformID: entity.PlatformID,
		ProjectID:  entity.ProjectID,
		Token:      entity.Token,
		EventType:  entity.EventType,
		CreatedAt:  entity.CreatedAt,
	}, nil
}
