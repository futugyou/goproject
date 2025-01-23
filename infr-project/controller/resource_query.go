package controller

import (
	"context"
	"net/http"
	"os"

	"github.com/futugyou/infr-project/application"
	"github.com/futugyou/infr-project/extensions"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ResourceQueryController struct {
}

func NewResourceQueryController() *ResourceQueryController {
	return &ResourceQueryController{}
}

func (c *ResourceQueryController) GetResource(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createResourceQueryService, func(ctx context.Context, service *application.ResourceQueryService, _ struct{}) (interface{}, error) {
		return service.GetResource(ctx, id)
	})
}

func (c *ResourceQueryController) GetAllResource(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createResourceQueryService, func(ctx context.Context, service *application.ResourceQueryService, _ struct{}) (interface{}, error) {
		return service.GetAllResources(ctx)
	})
}

func createResourceQueryService(ctx context.Context) (*application.ResourceQueryService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("query_db_name"),
		ConnectString: os.Getenv("query_mongodb_url"),
	}

	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	client, err := extensions.RedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}

	queryRepo := infra.NewResourceQueryRepository(mongoclient, config)

	unitOfWork, err := infra.NewMongoUnitOfWork(mongoclient)
	if err != nil {
		return nil, err
	}

	return application.NewResourceQueryService(queryRepo, client, unitOfWork), nil
}
