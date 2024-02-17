package cash

import (
	"github.com/futugyou/alphavantage-server/core"
)

type ICashRepository interface {
	core.IRepository[CashEntity, string]
}

type CashRepository struct {
	*core.MongoRepository[CashEntity, string]
}

func NewCashRepository(config core.DBConfig) *CashRepository {
	baseRepo := core.NewMongoRepository[CashEntity, string](config)
	return &CashRepository{baseRepo}
}
