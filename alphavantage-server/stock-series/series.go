package stockSeries

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage/enums"
)

func SyncStockSeriesData(symbol string) {
	log.Println("stock series data sync start.")
	// get sync month
	month := GetStaockMonth(symbol)

	log.Printf("start to get %s data, month %s \n", symbol, month)

	// get data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewTimeSeriesClient(apikey)
	p := alphavantage.TimeSeriesIntradayParameter{
		Symbol:   symbol,
		Interval: enums.T60min,
		Dictionary: map[string]string{
			"outputsize": "full",
			"month":      month,
		},
	}
	s, err := client.TimeSeriesIntraday(p)
	if err != nil {
		log.Println(err)
		return
	}

	// create insert data
	data := make([]StockSeriesEntity, 0)
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

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}
	repository := NewStockSeriesRepository(config)
	r, err := repository.InsertMany(context.Background(), data, StockFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	// update month
	if r.UpsertedCount > 0 || checkTime(month) {
		UpdateStaockMonth(month, symbol)
	}

	log.Println("stock series data sync finish")
}

func StockFilter(e StockSeriesEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}, {Key: "time", Value: e.Time}}
}

func checkTime(month string) bool {
	t, _ := time.Parse("2006-01", month)
	tt, _ := time.Parse("2006-01", time.Now().Format("2006-01"))
	return t.Before(tt)
}

func StockSeriesData() ([]StockSeriesEntity, error) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}
	repository := NewStockSeriesRepository(config)
	return repository.GetAll(context.Background())
}
