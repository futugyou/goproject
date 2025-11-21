package controller

import (
	"context"
	"net/http"

	"github.com/futugyou/resourcequeryservice/application"
	v1 "github.com/futugyou/resourcequeryservice/routes/v1"
)

type ResourceQueryController struct {
}

func NewResourceQueryController() *ResourceQueryController {
	return &ResourceQueryController{}
}

func (c *ResourceQueryController) GetResource(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateResourceQueryService, func(ctx context.Context, service *application.ResourceQueryService, _ struct{}) (any, error) {
		return service.GetResource(ctx, id)
	})
}

func (c *ResourceQueryController) GetAllResource(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateResourceQueryService, func(ctx context.Context, service *application.ResourceQueryService, _ struct{}) (any, error) {
		return service.GetAllResources(ctx)
	})
}
