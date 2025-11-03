package infrastructure

import (
	"context"
	"fmt"

	domaincore "github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/resourcequeryservice/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResourceQueryRepository struct {
	mongoimpl.BaseCRUD[domain.Resource]
}

func NewResourceQueryRepository(client *mongo.Client, config mongoimpl.DBConfig) *ResourceQueryRepository {
	if config.CollectionName == "" {
		config.CollectionName = "resources_query"
	}

	getID := func(a domain.Resource) string { return a.AggregateId() }

	return &ResourceQueryRepository{
		BaseCRUD: *mongoimpl.NewBaseCRUD(client, config, getID),
	}
}

func (r *ResourceQueryRepository) GetResourceByName(ctx context.Context, name string) (*domain.Resource, error) {
	var page, size int = 1, 1
	query := domaincore.NewQuery().
		Eq("name", name).
		Build()
	condition := domaincore.NewQueryOptions(&page, &size, nil, query)
	ent, err := r.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(ent) == 0 {
		return nil, fmt.Errorf("%s with name %s", domaincore.DATA_NOT_FOUND_MESSAGE, name)
	}
	return &ent[0], nil
}

func (r *ResourceQueryRepository) GetAllResource(ctx context.Context, page *int, size *int) ([]domain.Resource, error) {
	condition := domaincore.NewQueryOptions(page, size, nil, nil)
	return r.Find(ctx, condition)
}
