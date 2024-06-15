package apiadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/application"
	"github.com/futugyou/infr-project/extensions"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	models "github.com/futugyou/infr-project/view_models"
)

func CreateProject(w http.ResponseWriter, r *http.Request) {
	service, err := createProjectService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.CreateProject(aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func GetProject(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createProjectService()

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.GetProject(id)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func GetAllProject(w http.ResponseWriter, r *http.Request) {
	service, err := createProjectService()

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.GetAllProject()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func UpdateProject(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createProjectService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.UpdateProject(id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func UpdateProjectPlatform(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createProjectService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux []models.UpdateProjectPlatformRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.UpdateProjectPlatform(id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func UpdateProjectDesign(id string, w http.ResponseWriter, r *http.Request) {
	service, err := createProjectService()
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux []models.UpdateProjectDesignRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.UpdateProjectDesign(id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func createProjectService() (*application.ProjectService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	repo := infra.NewProjectRepository(client, config)
	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		return nil, err
	}

	return application.NewProjectService(unitOfWork, repo), nil
}
