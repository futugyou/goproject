package controller

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

type ProjectController struct {
}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

func (c *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createProjectService(ctx)
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

	res, err := service.CreateProject(ctx, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *ProjectController) GetProject(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createProjectService(ctx)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.GetProject(ctx, id)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *ProjectController) GetAllProject(w http.ResponseWriter, r *http.Request, page *int, size *int) {
	ctx := r.Context()
	service, err := createProjectService(ctx)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.GetAllProject(ctx, page, size)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *ProjectController) UpdateProject(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createProjectService(ctx)
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

	res, err := service.UpdateProject(ctx, id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *ProjectController) UpdateProjectPlatform(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createProjectService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux []models.UpdateProjectPlatformRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Var(aux, "required,gt=0,dive"); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.UpdateProjectPlatform(ctx, id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *ProjectController) UpdateProjectDesign(id string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createProjectService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux []models.UpdateProjectDesignRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Var(aux, "gt=0,dive"); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.UpdateProjectDesign(ctx, id, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func createProjectService(ctx context.Context) (*application.ProjectService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
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
