package api

import (
	"encoding/json"

	_ "github.com/joho/godotenv/autoload"

	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"net/http"

	"github.com/futugyou/infr-project/api/internal"
	"github.com/futugyou/infr-project/application"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/resource"
)

func ResourceDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if internal.CorsForVercel(w, r) {
		return
	}

	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	eventStore := infra.NewMongoEventStore(client, config, "resource_events", resource.CreateEvent)
	snapshotStore := infra.NewMongoSnapshotStore[*resource.Resource](client, config)
	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	resourceService := application.NewResourceService(eventStore, snapshotStore, unitOfWork)

	op := r.URL.Query().Get("optype")
	switch op {
	case "create":
		createResource(resourceService, r, w)
	case "get":
		getResource(resourceService, r, w)
	case "update":
		updateResource(resourceService, r, w)
	case "delete":
		deleteResource(resourceService, r, w)
	case "history":
		historyResource(resourceService, r, w)
	default:
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
	}
}

func historyResource(service *application.ResourceService, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	list, err := service.AllVersionResource(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	body, _ := json.Marshal(list)
	w.Write(body)
	w.WriteHeader(200)
}

func deleteResource(service *application.ResourceService, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	err := service.DeleteResource(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func updateResource(service *application.ResourceService, r *http.Request, w http.ResponseWriter) {
	aux := &struct {
		Id   string `json:"id"`
		Data string `json:"data"`
	}{}

	err := json.NewDecoder(r.Body).Decode(aux)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	err = service.UpdateResourceDate(aux.Id, aux.Data)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
	}
}

func getResource(service *application.ResourceService, r *http.Request, w http.ResponseWriter) {
	id := r.URL.Query().Get("id")
	res, err := service.CurrentResource(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, _ := json.Marshal(res)
	w.Write(body)
	w.WriteHeader(200)
}

func createResource(service *application.ResourceService, r *http.Request, w http.ResponseWriter) {
	aux := &struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Data string `json:"data"`
	}{}

	err := json.NewDecoder(r.Body).Decode(aux)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	resourceType := resource.GetResourceType(aux.Type)
	res, err := service.CreateResource(aux.Name, resourceType, aux.Data)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
	}

	body, _ := json.Marshal(res)
	w.Write(body)
	w.WriteHeader(200)
}
