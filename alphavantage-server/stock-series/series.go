package stockSeries

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage/enums"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SyncStockSeriesData(symbol string) {
	log.Println("stock series data sync start.")

	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	month := "2000-01"
	repo := NewStockSeriesConfigRepository(config)
	c, _ := repo.GetAll(context.Background())
	if len(c) > 0 && len(c[0].Month) > 0 {
		month = c[0].Month
	}

	// create series list
	data := make([]StockSeriesEntity, 0)
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewTimeSeriesClient(apikey)

	log.Printf("start to get %s data \n", symbol)
	p := alphavantage.TimeSeriesIntradayParameter{
		Symbol:   symbol,
		Interval: enums.T60min,
	}
	dic := make(map[string]string)
	dic["outputsize"] = "full"
	dic["month"] = month
	p.Dictionary = dic
	s, err := client.TimeSeriesIntraday(p)
	if err != nil {
		log.Println(err)
		return
	}
	for ii := 0; ii < len(s); ii++ {
		e := StockSeriesEntity{
			Id:     s[ii].Symbol + s[ii].Time.Format("2006-01-02 15:04:05"),
			Symbol: s[ii].Symbol,
			Time:   s[ii].Time,
			Open:   s[ii].Open,
			High:   s[ii].High,
			Low:    s[ii].Low,
			Close:  s[ii].Close,
			Volume: s[ii].Volume,
		}
		data = append(data, e)
	}

	repository := NewStockSeriesRepository(config)
	repository.InsertMany(context.Background(), data, StockFilter)

	t, _ := time.Parse("2006-01", month)
	configs := []StockSeriesConfigEntity{
		{
			Month:  t.AddDate(0, 1, 0).Format("2006-01"),
			Filter: "month",
		},
	}
	repo.InsertMany(context.Background(), configs, StockConfigFilter)
	log.Println("stock series data sync finish")
}

func StockFilter(e StockSeriesEntity) primitive.D {
	return bson.D{{Key: "symbol", Value: e.Symbol}, {Key: "time", Value: e.Time}}
}

func StockConfigFilter(e StockSeriesConfigEntity) primitive.D {
	return bson.D{{Key: "filter", Value: e.Filter}}
}
