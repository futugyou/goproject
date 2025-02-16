package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/application"
	"github.com/futugyou/infr-project/command"
	"github.com/futugyou/infr-project/extensions"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	publisher "github.com/futugyou/infr-project/infrastructure_qstash"
	models "github.com/futugyou/infr-project/view_models"
)

type PlatformController struct {
}

func NewPlatformController() *PlatformController {
	return &PlatformController{}
}

func (c *PlatformController) DeletePlatformProject(idOrName string, projectId string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (interface{}, error) {
		return service.DeleteProject(ctx, idOrName, projectId)
	})
}

func (c *PlatformController) UpsertPlatformProject(idOrName string, projectId string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, req models.UpdatePlatformProjectRequest) (interface{}, error) {
		return service.UpsertProject(ctx, idOrName, projectId, req)
	})
}

func (c *PlatformController) CreatePlatform(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, req models.CreatePlatformRequest) (interface{}, error) {
		return service.CreatePlatform(ctx, req)
	})
}

func (c *PlatformController) GetPlatform(idOrName string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (interface{}, error) {
		return service.GetPlatform(ctx, idOrName)
	})
}

func (c *PlatformController) GetPlatformProject(idOrName string, projectId string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (interface{}, error) {
		return service.GetPlatformProject(ctx, idOrName, projectId)
	})
}

func (c *PlatformController) SearchPlatforms(w http.ResponseWriter, r *http.Request, request models.SearchPlatformsRequest) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (interface{}, error) {
		return service.SearchPlatforms(ctx, request)
	})
}

func (c *PlatformController) UpdatePlatformHook(idOrName string, projectId string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, req models.UpdatePlatformWebhookRequest) (interface{}, error) {
		return service.UpsertWebhook(ctx, idOrName, projectId, req)
	})
}

func (c *PlatformController) DeletePlatformHook(request models.RemoveWebhookRequest, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (interface{}, error) {
		return service.RemoveWebhook(ctx, request)
	})
}

func (c *PlatformController) UpdatePlatform(idOrName string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, req models.UpdatePlatformRequest) (interface{}, error) {
		return service.UpdatePlatform(ctx, idOrName, req)
	})
}

func (c *PlatformController) DeletePlatform(idOrName string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (interface{}, error) {
		return service.DeletePlatform(ctx, idOrName)
	})
}

func (c *PlatformController) RecoveryPlatform(idOrName string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createPlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (interface{}, error) {
		return service.RecoveryPlatform(ctx, idOrName)
	})
}

func createPlatformService(ctx context.Context) (*application.PlatformService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		return nil, err
	}

	redisClient, err := extensions.RedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}

	repo := infra.NewPlatformRepository(client, config)
	vaultRepo := infra.NewVaultRepository(client, config)
	eventPublisher := publisher.NewQStashEventPulisher(os.Getenv("QSTASH_TOKEN"), os.Getenv("QSTASH_DESTINATION"))
	vaultService := application.NewVaultService(unitOfWork, vaultRepo, eventPublisher)
	return application.NewPlatformService(unitOfWork, repo, vaultService, redisClient, eventPublisher), nil
}

func (c *PlatformController) CreatePlatformV2(cqrsRoute *command.Router, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var aux command.CreatePlatformCommand
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	commandBus := cqrsRoute.CommandBus
	//TODO: this err is not Handle's response, i dot know what it is
	if err := commandBus.Send(ctx, aux); err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, nil, 200)
}
