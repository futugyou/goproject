package controller

import (
	"encoding/json"
	"fmt"

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

func (c *Controller) GetResourceHistory(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createResourceService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.AllVersionResource(id)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *Controller) DeleteResource(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createResourceService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	err = service.DeleteResource(id)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, "ok", 200)
}

func (c *Controller) UpdateResource(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createResourceService()
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

	err = service.UpdateResource(id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, "ok", 200)
}

func (c *Controller) CreateResource(w http.ResponseWriter, r *http.Request) {
	service, err := createResourceService()
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

	res, err := service.CreateResource(aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *Controller) GetResource(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createResourceQueryService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.CurrentResource(id)
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

func (c *Controller) GetAllResource(w http.ResponseWriter, r *http.Request) {
	service, err := createResourceQueryService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.GetAllResources()
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

func createResourceService() (*application.ResourceService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	eventStore := infra.NewMongoEventStore(client, config, "resource_events", resource.CreateEvent)
	snapshotStore := infra.NewMongoSnapshotStore[*resource.Resource](client, config)
	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		return nil, err
	}

	queryRepo, err := createResourceQueryRepository()
	if err != nil {
		return nil, err
	}

	return application.NewResourceService(eventStore, snapshotStore, unitOfWork, queryRepo), nil
}

func createResourceQueryService() (*application.ResourceQueryService, error) {
	queryRepo, err := createResourceQueryRepository()
	if err != nil {
		return nil, err
	}
	return application.NewResourceQueryService(queryRepo), nil
}

func createResourceQueryRepository() (*infra.ResourceQueryRepository, error) {
	config := infra.QueryDBConfig{
		DBName:        os.Getenv("query_db_name"),
		ConnectString: os.Getenv("query_mongodb_url"),
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	return infra.NewResourceQueryRepository(client, config), nil
}
