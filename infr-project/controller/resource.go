package controller

import (
	_ "github.com/joho/godotenv/autoload"

	"context"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/qstash"

	"github.com/futugyou/infr-project/application"
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
	handleRequest(w, r, createResourceService, func(ctx context.Context, service *application.ResourceService, _ struct{}) (interface{}, error) {
		return service.AllVersionResource(ctx, id)
	})
}

func (c *ResourceController) DeleteResource(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createResourceService, func(ctx context.Context, service *application.ResourceService, _ struct{}) (interface{}, error) {
		return "ok", service.DeleteResource(ctx, id)
	})
}

func (c *ResourceController) UpdateResource(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createResourceService, func(ctx context.Context, service *application.ResourceService, req models.UpdateResourceRequest) (interface{}, error) {
		return "ok", service.UpdateResource(ctx, id, req)
	})
}

func (c *ResourceController) CreateResource(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createResourceService, func(ctx context.Context, service *application.ResourceService, req models.CreateResourceRequest) (interface{}, error) {
		return service.CreateResource(ctx, req)
	})
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

	qstashClient := qstash.NewClient(os.Getenv("QSTASH_TOKEN"))
	return application.NewResourceService(eventStore, snapshotStore, unitOfWork, qstashClient), nil
}
