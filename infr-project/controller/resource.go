package controller

import (
	"encoding/json"

	_ "github.com/joho/godotenv/autoload"

	"context"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/application"
	"github.com/futugyou/infr-project/extensions"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/resource"
	models "github.com/futugyou/infr-project/view_models"
)

type ResourceController struct {
}

func NewResourceController() *ResourceController {
	return &ResourceController{}
}

func (c *ResourceController) GetResourceHistory(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createResourceService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.AllVersionResource(ctx, id)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *ResourceController) DeleteResource(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createResourceService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	err = service.DeleteResource(ctx, id)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, "ok", 200)
}

func (c *ResourceController) UpdateResource(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createResourceService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.UpdateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	err = service.UpdateResource(ctx, id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, "ok", 200)
}

func (c *ResourceController) CreateResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createResourceService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.CreateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.CreateResource(ctx, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func createResourceService(ctx context.Context) (*application.ResourceService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	eventStore := infra.NewMongoEventStore(client, config, "resource_events", resource.CreateEvent)
	snapshotStore := infra.NewMongoSnapshotStore[*resource.Resource](client, config)
	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		return nil, err
	}

	return application.NewResourceService(eventStore, snapshotStore, unitOfWork), nil
}
