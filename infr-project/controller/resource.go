package controller

import (
	_ "github.com/joho/godotenv/autoload"

	"context"
	"net/http"

	"github.com/futugyou/resourceservice/application"
	v1 "github.com/futugyou/resourceservice/routes/v1"
	"github.com/futugyou/resourceservice/viewmodel"
)

type ResourceController struct {
}

func NewResourceController() *ResourceController {
	return &ResourceController{}
}

func (c *ResourceController) GetResourceHistory(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateResourceService, func(ctx context.Context, service *application.ResourceService, _ struct{}) (interface{}, error) {
		return service.AllVersionResource(ctx, id)
	})
}

func (c *ResourceController) DeleteResource(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateResourceService, func(ctx context.Context, service *application.ResourceService, _ struct{}) (interface{}, error) {
		return nil, service.DeleteResource(ctx, id)
	})
}

func (c *ResourceController) UpdateResource(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateResourceService, func(ctx context.Context, service *application.ResourceService, req viewmodel.UpdateResourceRequest) (interface{}, error) {
		return nil, service.UpdateResource(ctx, id, req)
	})
}

func (c *ResourceController) CreateResource(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateResourceService, func(ctx context.Context, service *application.ResourceService, req viewmodel.CreateResourceRequest) (interface{}, error) {
		return service.CreateResource(ctx, req)
	})
}
