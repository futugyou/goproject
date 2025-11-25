package controller

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/futugyou/projectservice/application"
	v1 "github.com/futugyou/projectservice/routes/v1"
	"github.com/futugyou/projectservice/viewmodel"
)

type ProjectController struct {
}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

func (c *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateProjectService, func(ctx context.Context, service *application.ProjectService, req viewmodel.CreateProjectRequest) (any, error) {
		return service.CreateProject(ctx, req)
	})
}

func (c *ProjectController) GetProject(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateProjectService, func(ctx context.Context, service *application.ProjectService, _ struct{}) (any, error) {
		return service.GetProject(ctx, id)
	})
}

func (c *ProjectController) GetAllProject(w http.ResponseWriter, r *http.Request, page *int, size *int) {
	handleRequest(w, r, v1.CreateProjectService, func(ctx context.Context, service *application.ProjectService, _ struct{}) (any, error) {
		return service.GetAllProject(ctx, page, size)
	})
}

func (c *ProjectController) UpdateProject(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateProjectService, func(ctx context.Context, service *application.ProjectService, req viewmodel.UpdateProjectRequest) (any, error) {
		return nil, service.UpdateProject(ctx, id, req)
	})
}

func (c *ProjectController) UpdateProjectPlatform(id string, w http.ResponseWriter, r *http.Request) {
	handleRequestUseSpecValidate(w, r, v1.CreateProjectService,
		func(v *validator.Validate, req *[]viewmodel.UpdateProjectPlatformRequest) error {
			return v.Var(*req, "required,gt=0,dive")
		},
		func(ctx context.Context, service *application.ProjectService, req []viewmodel.UpdateProjectPlatformRequest) (any, error) {
			return nil, service.UpdateProjectPlatform(ctx, id, req)
		},
	)
}

func (c *ProjectController) UpdateProjectDesign(id string, w http.ResponseWriter, r *http.Request) {
	handleRequestUseSpecValidate(w, r, v1.CreateProjectService,
		func(v *validator.Validate, req *[]viewmodel.UpdateProjectDesignRequest) error {
			return v.Var(*req, "gt=0,dive")
		},
		func(ctx context.Context, service *application.ProjectService, req []viewmodel.UpdateProjectDesignRequest) (any, error) {
			return nil, service.UpdateProjectDesign(ctx, id, req)
		},
	)
}
