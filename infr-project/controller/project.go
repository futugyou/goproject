package controller

import (
	"context"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/application"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	models "github.com/futugyou/infr-project/view_models"
	"github.com/go-playground/validator/v10"
)

type ProjectController struct {
}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

func (c *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createProjectService, func(ctx context.Context, service *application.ProjectService, req models.CreateProjectRequest) (interface{}, error) {
		return service.CreateProject(ctx, req)
	})
}

func (c *ProjectController) GetProject(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createProjectService, func(ctx context.Context, service *application.ProjectService, _ struct{}) (interface{}, error) {
		return service.GetProject(ctx, id)
	})
}

func (c *ProjectController) GetAllProject(w http.ResponseWriter, r *http.Request, page *int, size *int) {
	handleRequest(w, r, createProjectService, func(ctx context.Context, service *application.ProjectService, _ struct{}) (interface{}, error) {
		return service.GetAllProject(ctx, page, size)
	})
}

func (c *ProjectController) UpdateProject(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createProjectService, func(ctx context.Context, service *application.ProjectService, req models.UpdateProjectRequest) (interface{}, error) {
		return service.UpdateProject(ctx, id, req)
	})
}

func (c *ProjectController) UpdateProjectPlatform(id string, w http.ResponseWriter, r *http.Request) {
	handleRequestUseSpecValidate(w, r, createProjectService,
		func(v *validator.Validate, req *[]models.UpdateProjectPlatformRequest) error {
			return v.Var(*req, "required,gt=0,dive")
		},
		func(ctx context.Context, service *application.ProjectService, req []models.UpdateProjectPlatformRequest) (interface{}, error) {
			return service.UpdateProjectPlatform(ctx, id, req)
		},
	)
}

func (c *ProjectController) UpdateProjectDesign(id string, w http.ResponseWriter, r *http.Request) {
	handleRequestUseSpecValidate(w, r, createProjectService,
		func(v *validator.Validate, req *[]models.UpdateProjectDesignRequest) error {
			return v.Var(*req, "gt=0,dive")
		},
		func(ctx context.Context, service *application.ProjectService, req []models.UpdateProjectDesignRequest) (interface{}, error) {
			return service.UpdateProjectDesign(ctx, id, req)
		},
	)
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
