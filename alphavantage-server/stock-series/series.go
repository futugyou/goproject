package stockSeries

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
)

func SyncStockSeriesData(ctx context.Context, symbol string) bool {
	log.Println("stock series data sync start.")
	// get sync month
	month := GetStockMonth(ctx, symbol)

	log.Printf("start to get %s data, month %s \n", symbol, month)

	// get data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	// The intraday API is now a Premium API, so switched to the free daily API.
	client := alphavantage.NewTimeSeriesClient(apikey)
	p := alphavantage.TimeSeriesDailyParameter{
		Symbol: symbol,
		// "outputsize": "full", is Premium option, month is no need, so remove it
		// Dictionary: map[string]string{
		// 	"outputsize": "full",
		// 	"month":      month,
		// },
	}
	s, err := client.TimeSeriesDaily(p)
	// alphavantage will throw 'Invalid API call' when no data, there is no way to distinguish 'no data' error from other errors.
	if err != nil {
		log.Println(err)
		// Stop outside loop when 'API rate limit is 25 requests per day'
		// Donot stop outside loop when other error
		if strings.Contains(err.Error(), "Thank you for using Alpha Vantage") {
			return true
		}
	}

	if len(s) > 0 {
		// create insert data
		data := make([]StockSeriesEntity, 0)
		for ii := range s {
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
		r, err := repository.InsertMany(ctx, data, StockFilter)
		if err != nil {
			log.Println(err)
			return false
		}

		r.String()
	}

	// update month
	needUpdate := checkTime(month)
	if needUpdate {
		UpdateStockMonth(ctx, month, symbol)
	}

	log.Println("stock series data sync finish")
	return !needUpdate
}

func StockFilter(e StockSeriesEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}, {Key: "time", Value: e.Time}}
}

func checkTime(month string) bool {
	t, _ := time.Parse("2006-01", month)
	tt, _ := time.Parse("2006-01", time.Now().Format("2006-01"))
	return t.Before(tt)
}

func StockSeriesData(ctx context.Context, symbol string, year string) ([]StockSeriesEntity, error) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewStockSeriesRepository(config)
	start, _ := time.Parse("2006", year)
	end := start.AddDate(1, 0, 0)
	return repository.GetWithFilter(ctx, []core.DataFilterItem{{
		Key:   "symbol",
		Value: symbol,
	}, {
		Key:   "time",
		Value: map[string]any{"$gte": start, "$lt": end},
	}})
}
