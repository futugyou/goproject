package stockSeries

import (
	"context"
	"os"
	"time"

	"github.com/futugyou/alphavantage-server/core"
)

func GetStaockMonth(ctx context.Context, symbol string) string {
	month := "2000-01"

	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repo := NewStockSeriesConfigRepository(config)
	conf, _ := repo.GetOne(ctx, []core.DataFilterItem{{Key: "symbol", Value: symbol}})
	if conf != nil {
		month = conf.Month
	}

	return month
}

func UpdateStaockMonth(ctx context.Context, month string, symbol string) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	t, _ := time.Parse("2006-01", month)
	repo := NewStockSeriesConfigRepository(config)
	repo.Update(ctx,
		StockSeriesConfigEntity{
			Month:  t.AddDate(0, 1, 0).Format("2006-01"),
			Symbol: symbol,
		},
		[]core.DataFilterItem{{Key: "symbol", Value: symbol}},
	)
}

func StockConfigFilter(e StockSeriesConfigEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}}
}
