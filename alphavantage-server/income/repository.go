package income

import (
	"github.com/futugyou/alphavantage-server/core"
)

type IIncomeRepository interface {
	core.IRepository[IncomeEntity, string]
}

type IncomeRepository struct {
	*core.MongoRepository[IncomeEntity, string]
}

func NewIncomeRepository(config core.DBConfig) *IncomeRepository {
	baseRepo := core.NewMongoRepository[IncomeEntity, string](config)
	return &IncomeRepository{baseRepo}
}
