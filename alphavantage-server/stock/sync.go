package stock

import (
	"context"
	"os"

	"log"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
)

func SyncStockSymbolData(symbol string) {
	log.Printf("%s company data sync start. \n", symbol)
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	// get data from db
	repository := NewStockRepository(config)
	stock, _ := repository.GetOne(context.Background(), []core.DataFilterItem{{Key: "symbol", Value: symbol}})
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	if stock == nil {
		log.Printf("no org %s company data \n", symbol)
		// get data from SymbolSearch API
		client := alphavantage.NewTimeSeriesClient(apikey)
		pp := alphavantage.SymbolSearchParameter{
			Keywords: symbol,
		}

		symbolList, err := client.SymbolSearch(pp)
		if err != nil {
			log.Println(err)
			return
		}
		// only use one data
		for i := 0; i < 1; i++ {
			stock = &StockEntity{
				Id:          symbolList[i].Symbol,
				Symbol:      symbolList[i].Symbol,
				Name:        symbolList[i].Name,
				Type:        symbolList[i].Type,
				Region:      symbolList[i].Region,
				MarketOpen:  symbolList[i].MarketOpen,
				MarketClose: symbolList[i].MarketClose,
				Timezone:    symbolList[i].Timezone,
				Currency:    symbolList[i].Currency,
				MatchScore:  symbolList[i].MatchScore,
			}
		}
	}

	if stock == nil {
		stock = &StockEntity{}
	}

	// get data from CompanyOverview API
	client := alphavantage.NewFundamentalsClient(apikey)
	pp := alphavantage.CompanyOverviewParameter{
		Symbol: symbol,
	}

	result, err := client.CompanyOverview(pp)
	if err != nil {
		log.Println(err)
		return
	}

	stock.Symbol = result.Symbol
	stock.AssetType = result.AssetType
	stock.Name = result.Name
	stock.Description = result.Description
	stock.Cik = result.Cik
	stock.Exchange = result.Exchange
	stock.Currency = result.Currency
	stock.Country = result.Country
	stock.Sector = result.Sector
	stock.Industry = result.Industry
	stock.Address = result.Address
	stock.FiscalYearEnd = result.FiscalYearEnd
	stock.LatestQuarter = result.LatestQuarter
	stock.MarketCapitalization = result.MarketCapitalization
	stock.Ebitda = result.Ebitda
	stock.PERatio = result.PERatio
	stock.PEGRatio = result.PEGRatio
	stock.BookValue = result.BookValue
	stock.DividendPerShare = result.DividendPerShare
	stock.DividendYield = result.DividendYield
	stock.Eps = result.Eps
	stock.RevenuePerShareTTM = result.RevenuePerShareTTM
	stock.ProfitMargin = result.ProfitMargin
	stock.OperatingMarginTTM = result.OperatingMarginTTM
	stock.ReturnOnAssetsTTM = result.ReturnOnAssetsTTM
	stock.ReturnOnEquityTTM = result.ReturnOnEquityTTM
	stock.RevenueTTM = result.RevenueTTM
	stock.GrossProfitTTM = result.GrossProfitTTM
	stock.DilutedEPSTTM = result.DilutedEPSTTM
	stock.QuarterlyEarningsGrowthYOY = result.QuarterlyEarningsGrowthYOY
	stock.QuarterlyRevenueGrowthYOY = result.QuarterlyRevenueGrowthYOY
	stock.AnalystTargetPrice = result.AnalystTargetPrice
	stock.TrailingPE = result.TrailingPE
	stock.ForwardPE = result.ForwardPE
	stock.PriceToSalesRatioTTM = result.PriceToSalesRatioTTM
	stock.PriceToBookRatio = result.PriceToBookRatio
	stock.EVToRevenue = result.EVToRevenue
	stock.EVToEBITDA = result.EVToEBITDA
	stock.Beta = result.Beta
	stock.The52WeekHigh = result.The52WeekHigh
	stock.The52WeekLow = result.The52WeekLow
	stock.The50DayMovingAverage = result.The50DayMovingAverage
	stock.The200DayMovingAverage = result.The200DayMovingAverage
	stock.SharesOutstanding = result.SharesOutstanding
	stock.DividendDate = result.DividendDate
	stock.ExDividendDate = result.ExDividendDate

	repository.Update(context.Background(), *stock, []core.DataFilterItem{{Key: "symbol", Value: symbol}})

	log.Printf("%s company data sync end. \n", symbol)
}

func StockSymbolFilter(e StockEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}}
}

func StockSymbolDatas() ([]StockEntity, error) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewStockRepository(config)
	return repository.GetAll(context.Background())
}
