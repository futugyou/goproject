package earnings

import (
	"github.com/futugyou/alphavantage-server/core"
)

type IEarningsRepository interface {
	core.IRepository[EarningsEntity, string]
}

type EarningsRepository struct {
	*core.MongoRepository[EarningsEntity, string]
}

func NewEarningsRepository(config core.DBConfig) *EarningsRepository {
	baseRepo := core.NewMongoRepository[EarningsEntity, string](config)
	return &EarningsRepository{baseRepo}
}
