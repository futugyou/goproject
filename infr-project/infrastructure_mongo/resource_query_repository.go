package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/extensions"
	resourcequery "github.com/futugyou/infr-project/resource_query"
)

type ResourceQueryRepository struct {
	BaseRepository[resourcequery.Resource]
}

func NewResourceQueryRepository(client *mongo.Client, config DBConfig) *ResourceQueryRepository {
	return &ResourceQueryRepository{
		BaseRepository: *NewBaseRepository[resourcequery.Resource](client, config),
	}
}

func (r *ResourceQueryRepository) GetResourceByName(ctx context.Context, name string) (*resourcequery.Resource, error) {
	var page, size int = 1, 1
	condition := extensions.NewSearch(&page, &size, nil, map[string]interface{}{"name": name})
	ent, err := r.BaseRepository.GetWithCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(ent) == 0 {
		return nil, fmt.Errorf("%s with name %s", extensions.Data_Not_Found_Message, name)
	}
	return &ent[0], nil
}

func (r *ResourceQueryRepository) GetAllResource(ctx context.Context, page *int, size *int) ([]resourcequery.Resource, error) {
	condition := extensions.NewSearch(page, size, nil, nil)
	return r.BaseRepository.GetWithCondition(ctx, condition)
}
