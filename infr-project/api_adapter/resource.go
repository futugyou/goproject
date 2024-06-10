package apiadapter

import (
	"encoding/json"

	_ "github.com/joho/godotenv/autoload"

	"context"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/application"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/resource"
	models "github.com/futugyou/infr-project/view_models"
)

func GetResourceHistory(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createResourceService()

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	res, err := service.AllVersionResource(id)

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	body, _ := json.Marshal(res)
	w.Write(body)
	w.WriteHeader(200)
}

func DeleteResource(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createResourceService()

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	err = service.DeleteResource(id)

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	w.Write([]byte("ok"))
	w.WriteHeader(200)
}

func UpdateResource(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createResourceService()

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	var aux models.UpdateResourceRequest

	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&aux); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	err = service.UpdateResourceDate(id, aux.Data)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	w.Write([]byte("ok"))
	w.WriteHeader(200)
}

func CreateResource(w http.ResponseWriter, r *http.Request) {
	service, err := createResourceService()

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	var aux models.CreateResourceRequest

	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&aux); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	resourceType := resource.GetResourceType(aux.Type)
	res, err := service.CreateResource(aux.Name, resourceType, aux.Data)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	body, _ := json.Marshal(res)
	w.Write(body)
	w.WriteHeader(200)
}

func GetResource(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createResourceQueryService()

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	res, err := service.CurrentResource(id)

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	if res == nil || res.Id == "" {
		w.Write([]byte("resource not found"))
		w.WriteHeader(400)
		return
	}

	body, _ := json.Marshal(res)
	w.Write(body)
	w.WriteHeader(200)
}

func GetAllResource(w http.ResponseWriter, r *http.Request) {
	service, err := createResourceQueryService()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	res, err := service.GetAllResources()

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	if len(res) == 0 {
		w.Write([]byte("resource not found"))
		w.WriteHeader(400)
		return
	}

	body, _ := json.Marshal(res)
	w.Write(body)
	w.WriteHeader(200)
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

	return application.NewResourceService(eventStore, snapshotStore, unitOfWork), nil
}

func createResourceQueryService() (*application.ResourceQueryService, error) {
	config := infra.QueryDBConfig{
		DBName:        os.Getenv("query_db_name"),
		ConnectString: os.Getenv("query_mongodb_url"),
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	repo := infra.NewResourceQueryRepository(client, config)

	return application.NewResourceQueryService(repo), nil
}
