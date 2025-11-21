package controller

import (
	"context"
	"net/http"

	"github.com/futugyou/platformservice/application"
	v1 "github.com/futugyou/platformservice/routes/v1"
	"github.com/futugyou/platformservice/viewmodel"
)

type PlatformController struct {
}

func NewPlatformController() *PlatformController {
	return &PlatformController{}
}

func (c *PlatformController) DeletePlatformProject(idOrName string, projectId string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (any, error) {
		return nil, service.DeleteProject(ctx, idOrName, projectId)
	})
}

func (c *PlatformController) UpsertPlatformProject(idOrName string, projectId string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, req viewmodel.UpdatePlatformProjectRequest) (any, error) {
		return nil, service.UpsertProject(ctx, idOrName, projectId, req)
	})
}

func (c *PlatformController) CreatePlatform(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, req viewmodel.CreatePlatformRequest) (any, error) {
		return nil, service.CreatePlatform(ctx, req)
	})
}

func (c *PlatformController) GetPlatform(idOrName string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (any, error) {
		return service.GetPlatform(ctx, idOrName)
	})
}

func (c *PlatformController) GetProviderProjectList(idOrName string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (any, error) {
		return service.GetProviderProjectList(ctx, idOrName)
	})
}

func (c *PlatformController) ImportProjectsFromProvider(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, req viewmodel.ImportProjectsRequest) (any, error) {
		return nil, service.ImportProjectsFromProvider(ctx, req.PlatformID, req.ProjectIDs)
	})
}

func (c *PlatformController) GetPlatformProject(idOrName string, projectId string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (any, error) {
		return service.GetPlatformProject(ctx, idOrName, projectId)
	})
}

func (c *PlatformController) SearchPlatforms(w http.ResponseWriter, r *http.Request, request viewmodel.SearchPlatformsRequest) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (any, error) {
		return service.SearchPlatforms(ctx, request)
	})
}

func (c *PlatformController) UpdatePlatform(idOrName string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, req viewmodel.UpdatePlatformRequest) (any, error) {
		return nil, service.UpdatePlatform(ctx, idOrName, req)
	})
}

func (c *PlatformController) DeletePlatform(idOrName string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (any, error) {
		return nil, service.DeletePlatform(ctx, idOrName)
	})
}

func (c *PlatformController) RecoveryPlatform(idOrName string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreatePlatformService, func(ctx context.Context, service *application.PlatformService, _ struct{}) (any, error) {
		return nil, service.RecoveryPlatform(ctx, idOrName)
	})
}
