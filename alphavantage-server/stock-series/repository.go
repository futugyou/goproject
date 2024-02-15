package stockSeries

import (
	"github.com/futugyou/alphavantage-server/core"
)

type IStockSeriesRepository interface {
	core.IRepository[StockSeriesEntity, string]
}

type StockSeriesRepository struct {
	*core.MongoRepository[StockSeriesEntity, string]
}

func NewStockSeriesRepository(config core.DBConfig) *StockSeriesRepository {
	baseRepo := core.NewMongoRepository[StockSeriesEntity, string](config)
	return &StockSeriesRepository{baseRepo}
}

type IStockSeriesConfigRepository interface {
	core.IRepository[StockSeriesConfigEntity, string]
}

type StockSeriesConfigRepository struct {
	*core.MongoRepository[StockSeriesConfigEntity, string]
}

func NewStockSeriesConfigRepository(config core.DBConfig) *StockSeriesConfigRepository {
	baseRepo := core.NewMongoRepository[StockSeriesConfigEntity, string](config)
	return &StockSeriesConfigRepository{baseRepo}
}
