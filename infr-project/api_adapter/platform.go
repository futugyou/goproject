package apiadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/application"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/platform"
	models "github.com/futugyou/infr-project/view_models"
)

func CreatePlatform(w http.ResponseWriter, r *http.Request) {
	service, err := createPlatformService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.CreatePlatformRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.CreatePlatform(aux.Name, aux.Url, aux.Rest, aux.Property)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func GetPlatform(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createPlatformService()

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.GetPlatform(id)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func GetAllPlatform(w http.ResponseWriter, r *http.Request) {
	service, err := createPlatformService()

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.GetAllPlatform()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func UpdatePlatformHook(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createPlatformService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux platform.Webhook
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.AddWebhook(id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func UpdatePlatform(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createPlatformService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.UpdatePlatformRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.UpdatePlatform(id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func DeletePlatform(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createPlatformService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	err = service.DeletePlatform(id)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, "ok", 200)
}

func createPlatformService() (*application.PlatformService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	repo := infra.NewPlatformRepository(client, config)
	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		return nil, err
	}

	return application.NewPlatformService(unitOfWork, repo), nil
}