package infrastructure

import (
	"context"
	"errors"
	"time"

	coreinfr "github.com/futugyou/domaincore/infrastructure"
	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/platformservice/application/service"
	"github.com/futugyou/platformservice/infrastructure/internal/dao"
	"github.com/futugyou/platformservice/infrastructure/internal/entity"
	"github.com/futugyou/platformservice/util"

	"go.mongodb.org/mongo-driver/mongo"
)

type EventHandler struct {
	log *dao.EventLogDao
	coreinfr.EventDispatcher
}

func NewEventHandler(client *mongo.Client, config mongoimpl.DBConfig, eventDispatcher coreinfr.EventDispatcher) *EventHandler {
	log := dao.NewEventLogDao(client, config)

	return &EventHandler{
		log:             log,
		EventDispatcher: eventDispatcher,
	}
}

func (e *EventHandler) DispatchIntegrationEvents(ctx context.Context, events []coreinfr.Event) error {
	// TODO: need token exchange
	token, ok := util.JWTFrom(ctx)
	if !ok || token == "" {
		return errors.New("missing jwt token in context")
	}

	entities := make([]entity.EventLogEntity, 0, len(events))
	for _, event := range events {
		switch t := event.(type) {
		case *CreateProviderProjectTriggeredEvent:
			entities = append(entities, entity.EventLogEntity{
				ID:         t.ID,
				PlatformID: t.PlatformID,
				ProjectID:  t.ProjectID,
				Token:      token,
				EventType:  t.EventType(),
				CreatedAt:  time.Now(),
			})
		case *CreateProviderWebhookTriggeredEvent:
			entities = append(entities, entity.EventLogEntity{
				ID:         t.ID,
				PlatformID: t.PlatformID,
				ProjectID:  t.ProjectID,
				Token:      token,
				EventType:  t.EventType(),
				CreatedAt:  time.Now(),
			})
		default: //screenshot do not need to be recorded
		}
	}

	err := e.log.InsertMany(ctx, entities)
	if err != nil {
		return err
	}

	return e.EventDispatcher.DispatchIntegrationEvents(ctx, events)
}

// Get implements service.EventHandler.
func (e *EventHandler) Get(ctx context.Context, eventID string) (*service.EventLog, error) {
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
