package stockSeries

import (
	"context"
	"os"
	"time"

	"github.com/futugyou/alphavantage-server/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetStaockMonth() string {
	month := "2000-01"

	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repo := NewStockSeriesConfigRepository(config)
	c, _ := repo.GetAll(context.Background())
	if len(c) > 0 && len(c[0].Month) > 0 {
		month = c[0].Month
	}

	return month
}

func UpdateStaockMonth(month string) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	t, _ := time.Parse("2006-01", month)
	configs := []StockSeriesConfigEntity{
		{
			Month:  t.AddDate(0, 1, 0).Format("2006-01"),
			Filter: "month",
		},
	}
	repo := NewStockSeriesConfigRepository(config)
	repo.InsertMany(context.Background(), configs, StockConfigFilter)
}

func StockConfigFilter(e StockSeriesConfigEntity) primitive.D {
	return bson.D{{Key: "filter", Value: e.Filter}}
}
