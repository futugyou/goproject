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
	models "github.com/futugyou/infr-project/view_models"
)

type PlatformController struct {
}

func NewPlatformController() *PlatformController {
	return &PlatformController{}
}

func (c *PlatformController) DeletePlatformProject(id string, projectId string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.DeleteProject(ctx, id, projectId)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *PlatformController) UpsertPlatformProject(id string, projectId string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.UpdatePlatformProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.UpsertProject(ctx, id, projectId, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *PlatformController) CreatePlatform(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.CreatePlatformRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.CreatePlatform(ctx, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *PlatformController) GetPlatform(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.GetPlatform(ctx, id)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *PlatformController) GetPlatformProject(id string, projectId string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.GetPlatformProject(ctx, id, projectId)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *PlatformController) SearchPlatforms(w http.ResponseWriter, r *http.Request, request models.SearchPlatformsRequest) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.SearchPlatforms(ctx, request)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *PlatformController) UpdatePlatformHook(id string, projectId string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.UpdatePlatformWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.UpsertWebhook(ctx, id, projectId, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *PlatformController) DeletePlatformHook(request models.RemoveWebhookRequest, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.RemoveWebhook(ctx, request)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *PlatformController) UpdatePlatform(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.UpdatePlatformRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.UpdatePlatform(ctx, id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *PlatformController) DeletePlatform(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.DeletePlatform(ctx, id)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *PlatformController) RecoveryPlatform(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createPlatformService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.RecoveryPlatform(ctx, id)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
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

	repo := infra.NewPlatformRepository(client, config)
	vaultRepo := infra.NewVaultRepository(client, config)
	vaultService := application.NewVaultService(unitOfWork, vaultRepo)
	return application.NewPlatformService(unitOfWork, repo, vaultService), nil
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
