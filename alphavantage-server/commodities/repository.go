package commodities

import (
	"github.com/futugyou/alphavantage-server/core"
)

type ICommoditiesRepository interface {
	core.IRepository[CommoditiesEntity, string]
}

type CommoditiesRepository struct {
	*core.MongoRepository[CommoditiesEntity, string]
}

func NewCommoditiesRepository(config core.DBConfig) *CommoditiesRepository {
	baseRepo := core.NewMongoRepository[CommoditiesEntity, string](config)
	return &CommoditiesRepository{baseRepo}
}
