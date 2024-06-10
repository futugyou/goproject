package infrastructure_mongo

import (
	models "github.com/futugyou/infr-project/view_models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResourceQueryRepository struct {
	BaseQueryRepository[models.ResourceDetail]
}

func NewResourceQueryRepository(client *mongo.Client, config QueryDBConfig) *ResourceQueryRepository {
	return &ResourceQueryRepository{
		BaseQueryRepository: *NewBaseQueryRepository[models.ResourceDetail](client, config),
	}
}
