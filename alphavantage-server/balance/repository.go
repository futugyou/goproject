package balance

import (
	"github.com/futugyou/alphavantage-server/core"
)

type IBalanceRepository interface {
	core.IRepository[BalanceEntity, string]
}

type BalanceRepository struct {
	*core.MongoRepository[BalanceEntity, string]
}

func NewBalanceRepository(config core.DBConfig) *BalanceRepository {
	baseRepo := core.NewMongoRepository[BalanceEntity, string](config)
	return &BalanceRepository{baseRepo}
}
