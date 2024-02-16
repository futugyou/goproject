package stockSeries

import (
	"context"
	"os"
	"time"

	"github.com/futugyou/alphavantage-server/core"
)

func GetStaockMonth(symbol string) string {
	month := "2000-01"

	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repo := NewStockSeriesConfigRepository(config)
	configList, _ := repo.GetAll(context.Background())
	if len(configList) > 0 {
		for _, config := range configList {
			if config.Symbol == symbol {
				month = config.Month
			}
		}
	}

	return month
}

func UpdateStaockMonth(month string, symbol string) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	t, _ := time.Parse("2006-01", month)
	configs := []StockSeriesConfigEntity{
		{
			Month:  t.AddDate(0, 1, 0).Format("2006-01"),
			Symbol: symbol,
		},
	}
	repo := NewStockSeriesConfigRepository(config)
	repo.InsertMany(context.Background(), configs, StockConfigFilter)
}

func StockConfigFilter(e StockSeriesConfigEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}}
}
