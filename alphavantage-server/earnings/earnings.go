package earnings

import (
	"context"
	"log"
	"os"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
)

func SyncEarningsData(symbol string) {
	log.Printf("%s earnings sentiment data sync start. \n", symbol)
	// get earnings data from alphavantage
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewFundamentalsClient(apikey)
	p := alphavantage.EarningsParameter{
		Symbol: symbol,
	}
	s, err := client.Earnings(p)
	if err != nil || s == nil {
		log.Println(err)
		return
	}

	// create Earnings list
	data := make([]EarningsEntity, 0)
	for ii := 0; ii < len(s.AnnualEarnings); ii++ {
		e := createEarningsEntity(s.AnnualEarnings[ii], symbol, "Annual")
		data = append(data, e)
	}
	for ii := 0; ii < len(s.QuarterlyEarnings); ii++ {
		e := createEarningsEntity(s.QuarterlyEarnings[ii], symbol, "Quarterly")
		data = append(data, e)
	}

	// insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewEarningsRepository(config)
	r, err := repository.InsertMany(context.Background(), data, EarningsFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("Earnings sentiment data sync finish")
}

func createEarningsEntity(earnings alphavantage.EarningsReport, symbol string, dataType string) EarningsEntity {
	return EarningsEntity{
		Id:                 symbol + "-" + dataType + "-" + earnings.FiscalDateEnding,
		Symbol:             symbol,
		DataType:           dataType,
		FiscalDateEnding:   earnings.FiscalDateEnding,
		ReportedDate:       earnings.ReportedDate,
		ReportedEPS:        earnings.ReportedEPS,
		EstimatedEPS:       earnings.EstimatedEPS,
		Surprise:           earnings.Surprise,
		SurprisePercentage: earnings.SurprisePercentage,
	}
}

func EarningsFilter(e EarningsEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}, {Key: "dataType", Value: e.DataType}, {Key: "fiscalDateEnding", Value: e.FiscalDateEnding}}
}
