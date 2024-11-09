package controller

import (
	"context"
	"fmt"
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
	ctx := r.Context()
	service, err := createResourceQueryService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.CurrentResource(ctx, id)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	if res == nil || res.Id == "" {
		handleError(w, fmt.Errorf("resource not found"), 400)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *ResourceQueryController) GetAllResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createResourceQueryService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.GetAllResources(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	if len(res) == 0 {
		handleError(w, fmt.Errorf("resource not found"), 400)
		return
	}

	writeJSONResponse(w, res, 200)
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
