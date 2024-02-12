package stock

import (
	"github.com/futugyou/alphavantage-server/core"
)

type IStockRepository interface {
	core.IRepository[StockEntity, string]
}

type StockRepository struct {
	*core.MongoRepository[StockEntity, string]
}

func NewStockRepository(config core.DBConfig) *StockRepository {
	baseRepo := core.NewMongoRepository[StockEntity, string](config)
	return &StockRepository{baseRepo}
}
