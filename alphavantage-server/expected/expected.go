package expected

import (
	"context"
	"log"
	"os"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage/enums"
)

func SyncExpectedData(ctx context.Context, symbol string) {
	log.Printf("%s expected sentiment data sync start. \n", symbol)
	// get expected data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewFundamentalsClient(apikey)
	p := alphavantage.EarningsCalendarParameter{
		Symbol:  symbol,
		Horizon: enums.H12month,
	}
	s, err := client.EarningsCalendar(p)
	if err != nil || s == nil {
		log.Println(err)
		return
	}

	// create Expected list
	data := make([]ExpectedEntity, 0)
	for _, v := range s {
		data = append(data, ExpectedEntity{
			Id:               v.Symbol + v.ReportDate.Format("2006-01-02"),
			Symbol:           symbol,
			Name:             v.Name,
			ReportDate:       v.ReportDate,
			FiscalDateEnding: v.FiscalDateEnding,
			Estimate:         v.Estimate,
			Currency:         v.Currency,
		})
	}

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewExpectedRepository(config)
	r, err := repository.InsertMany(ctx, data, ExpectedFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()

	log.Println("expected data sync finish")
}

func ExpectedFilter(e ExpectedEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}, {Key: "reportDate", Value: e.ReportDate}}
}
