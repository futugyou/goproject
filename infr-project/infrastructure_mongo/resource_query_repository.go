package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/extensions"
	models "github.com/futugyou/infr-project/view_models"
)

type ResourceQueryRepository struct {
	BaseQueryRepository[models.ResourceDetail]
}

func NewResourceQueryRepository(client *mongo.Client, config QueryDBConfig) *ResourceQueryRepository {
	return &ResourceQueryRepository{
		BaseQueryRepository: *NewBaseQueryRepository[models.ResourceDetail](client, config),
	}
}

func (r *ResourceQueryRepository) GetResourceByName(ctx context.Context, name string) (*models.ResourceDetail, error) {
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"name": name})
	ent, err := r.BaseQueryRepository.GetWithCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(ent) == 0 {
		return nil, fmt.Errorf("data not found with name %s", name)
	}
	return &ent[0], nil
}
