package controller

import (
	"context"
	"net/http"

	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/projectservice/application"
	"github.com/futugyou/projectservice/infrastructure"
	"github.com/futugyou/projectservice/options"
	"github.com/futugyou/projectservice/viewmodel"

	"github.com/go-playground/validator/v10"
)

type ProjectController struct {
}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

func (c *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createProjectService, func(ctx context.Context, service *application.ProjectService, req viewmodel.CreateProjectRequest) (interface{}, error) {
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
	handleRequest(w, r, createProjectService, func(ctx context.Context, service *application.ProjectService, req viewmodel.UpdateProjectRequest) (interface{}, error) {
		return nil, service.UpdateProject(ctx, id, req)
	})
}

func (c *ProjectController) UpdateProjectPlatform(id string, w http.ResponseWriter, r *http.Request) {
	handleRequestUseSpecValidate(w, r, createProjectService,
		func(v *validator.Validate, req *[]viewmodel.UpdateProjectPlatformRequest) error {
			return v.Var(*req, "required,gt=0,dive")
		},
		func(ctx context.Context, service *application.ProjectService, req []viewmodel.UpdateProjectPlatformRequest) (interface{}, error) {
			return nil, service.UpdateProjectPlatform(ctx, id, req)
		},
	)
}

func (c *ProjectController) UpdateProjectDesign(id string, w http.ResponseWriter, r *http.Request) {
	handleRequestUseSpecValidate(w, r, createProjectService,
		func(v *validator.Validate, req *[]viewmodel.UpdateProjectDesignRequest) error {
			return v.Var(*req, "gt=0,dive")
		},
		func(ctx context.Context, service *application.ProjectService, req []viewmodel.UpdateProjectDesignRequest) (interface{}, error) {
			return nil, service.UpdateProjectDesign(ctx, id, req)
		},
	)
}

func createProjectService(ctx context.Context) (*application.ProjectService, error) {
	option := options.New()
	mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.MongoDBURL)
	config := mongoimpl.DBConfig{
		DBName: option.DBName,
	}

	if err != nil {
		return nil, err
	}

	repo := infrastructure.NewProjectRepository(mongoclient, config)
	unitOfWork, err := mongoimpl.NewMongoUnitOfWork(mongoclient)
	if err != nil {
		return nil, err
	}

	return application.NewProjectService(unitOfWork, repo), nil
}
