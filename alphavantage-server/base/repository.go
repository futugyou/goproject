package base

import (
	"github.com/futugyou/alphavantage-server/core"
)

type IBaseDataRepository interface {
	core.IRepository[BaseDataEntity, string]
}

type BaseDataRepository struct {
	*core.MongoRepository[BaseDataEntity, string]
}

func NewBaseDataRepository(config core.DBConfig) *BaseDataRepository {
	baseRepo := core.NewMongoRepository[BaseDataEntity, string](config)
	return &BaseDataRepository{baseRepo}
}
