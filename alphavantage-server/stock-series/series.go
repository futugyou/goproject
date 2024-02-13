package stockSeries

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage-server/stock"
	"github.com/futugyou/alphavantage/enums"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// For example, month=2009-01. Any month in the last 20+ years since 2000-01 (January 2000) is supported.
// DOTO: add month table
func SyncStockSeriesData(month string) {
	log.Println("stock series data sync start.")
	// get stock symbol data from db
	list, err := stock.StockSymbolDatas()
	if err != nil {
		log.Println(err)
		return
	}

	if len(list) == 0 {
		return
	}

	// create series list
	data := make([]StockSeriesEntity, 0)
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewTimeSeriesClient(apikey)
	for i := 0; i < len(list); i++ {
		symbol := list[i].Symbol
		log.Printf("start to get %s data \n", symbol)
		p := alphavantage.TimeSeriesIntradayParameter{
			Symbol:   symbol,
			Interval: enums.T60min,
		}
		dic := make(map[string]string)
		dic["outputsize"] = "full"
		if len(strings.TrimSpace(month)) > 0 {
			dic["month"] = "month"
		}
		p.Dictionary = dic
		s, err := client.TimeSeriesIntraday(p)
		if err != nil {
			log.Println(err)
			continue
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
	}

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewStockSeriesRepository(config)
	repository.InsertMany(context.Background(), data, StockFilter)
	log.Println("stock series data sync finish")
}

func StockFilter(e StockSeriesEntity) primitive.D {
	return bson.D{{Key: "symbol", Value: e.Symbol}, {Key: "time", Value: e.Time}}
}
