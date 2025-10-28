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
	mongoimpl.BaseRepository[domain.Resource]
}

func NewResourceQueryRepository(client *mongo.Client, config mongoimpl.DBConfig) *ResourceQueryRepository {
	return &ResourceQueryRepository{
		BaseRepository: *mongoimpl.NewBaseRepository[domain.Resource](client, config),
	}
}

func (r *ResourceQueryRepository) GetResourceByName(ctx context.Context, name string) (*domain.Resource, error) {
	var page, size int = 1, 1
	condition := domaincore.NewQueryOptions(&page, &size, nil, map[string]any{"name": name})
	ent, err := r.BaseRepository.Find(ctx, condition)
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
	return r.BaseRepository.Find(ctx, condition)
}
