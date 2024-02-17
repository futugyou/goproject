package main

import (
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyou/alphavantage-server/balance"
	"github.com/futugyou/alphavantage-server/cash"
	"github.com/futugyou/alphavantage-server/commodities"
	"github.com/futugyou/alphavantage-server/earnings"
	"github.com/futugyou/alphavantage-server/expected"
	"github.com/futugyou/alphavantage-server/income"
	"github.com/futugyou/alphavantage-server/news"
	"github.com/futugyou/alphavantage-server/stock"
	stockSeries "github.com/futugyou/alphavantage-server/stock-series"
)

func main() {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		ProcessToRun()
		return
	}
	CommoditiesData()
}

func ProcessToRun() {
	w := time.Now().Weekday()
	if w == time.Saturday {
		CommoditiesData()
		// Economic Indicators data
	} else {
		index := int(time.Sunday)
		symbol := stock.StockList[index]
		SymbolData(symbol)

		Income(symbol)
		Balance(symbol)
		Cash(symbol)
		Earning(symbol)
		Expected(symbol)

		News(symbol)
		StockSeries(symbol)
	}

}

func SymbolData(symbol string) {
	stock.SyncStockSymbolData(symbol)
}

func StockSeries(symbol string) {
	stockSeries.SyncStockSeriesData(symbol)
}

func News(symbol string) {
	news.SyncNewsSentimentData(symbol)
}

func Income(symbol string) {
	income.SyncIncomeStatementData(symbol)
}

func Balance(symbol string) {
	balance.SyncBalanceSheetData(symbol)
}

func Cash(symbol string) {
	cash.SyncCashSheetData(symbol)
}

func Earning(symbol string) {
	earnings.SyncEarningsData(symbol)
}

func Expected(symbol string) {
	expected.SyncExpectedData(symbol)
}

func CommoditiesData() {
	commodities.SyncAllCommoditiesData()
}
