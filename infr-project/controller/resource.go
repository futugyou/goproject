package controller

import (
	_ "github.com/joho/godotenv/autoload"

	"context"
	"net/http"

	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/domaincore/qstashdispatcherimpl"

	"github.com/futugyou/resourceservice/application"
	"github.com/futugyou/resourceservice/infrastructure"
	"github.com/futugyou/resourceservice/options"

	"github.com/futugyou/resourceservice/viewmodel"
)

type ResourceController struct {
}

func NewResourceController() *ResourceController {
	return &ResourceController{}
}

func (c *ResourceController) GetResourceHistory(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createResourceService, func(ctx context.Context, service *application.ResourceService, _ struct{}) (interface{}, error) {
		return service.AllVersionResource(ctx, id)
	})
}

func (c *ResourceController) DeleteResource(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createResourceService, func(ctx context.Context, service *application.ResourceService, _ struct{}) (interface{}, error) {
		return nil, service.DeleteResource(ctx, id)
	})
}

func (c *ResourceController) UpdateResource(id string, w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createResourceService, func(ctx context.Context, service *application.ResourceService, req viewmodel.UpdateResourceRequest) (interface{}, error) {
		return nil, service.UpdateResource(ctx, id, req)
	})
}

func (c *ResourceController) CreateResource(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createResourceService, func(ctx context.Context, service *application.ResourceService, req viewmodel.CreateResourceRequest) (interface{}, error) {
		return service.CreateResource(ctx, req)
	})
}

func createResourceService(ctx context.Context) (*application.ResourceService, error) {
	option := options.New()
	mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.MongoDBURL)
	if err != nil {
		return nil, err
	}

	eventStore := infrastructure.NewResourceEventStore(mongoclient, option)
	snapshotStore := infrastructure.NewResourceSnapshotStore(mongoclient, option)

	unitOfWork, err := mongoimpl.NewMongoUnitOfWork(mongoclient)
	if err != nil {
		return nil, err
	}

	eventPublisher := qstashdispatcherimpl.NewQStashEventDispatcher(option.QstashToken, option.QstashDestination)
	return application.NewResourceService(eventStore, snapshotStore, unitOfWork, eventPublisher), nil
}
