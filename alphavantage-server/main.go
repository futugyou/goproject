package main

import (
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyou/alphavantage-server/balance"
	"github.com/futugyou/alphavantage-server/base"
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
	// commodities.CreateCommoditiesIndex()
}

func ProcessToRun() {
	d := time.Now().Day()
	w := time.Now().Weekday()
	count := 0
	// The number of available tokens is 25, but leave 3 as preparation
	totalCount := 22

	// 1. Runs on the 2nd of every month
	// This will consume 17 tokens, so return.
	if d == 2 {
		commodities.SyncMonthlyCommoditiesData()
		commodities.SyncMonthlyEconomicData()
		commodities.SyncQuarterlyEconomicData()
		commodities.SyncAnnualEconomicData()
		return
	}

	// Runs monthly stock sync from 3 to 3+len(StockList)
	// This will consume 5 tokens, so go on.
	for i := 0; i < len(stock.StockList); i++ {
		if d == i+3 {
			symbol := stock.StockList[i]
			income.SyncIncomeStatementData(symbol)
			balance.SyncBalanceSheetData(symbol)
			cash.SyncCashSheetData(symbol)
			earnings.SyncEarningsData(symbol)
			expected.SyncExpectedData(symbol)
			count += 5
		}
	}

	// Runs every Saturday or (3 of month and Sunday)
	// This will consume 5 tokens, so go on.
	if w == time.Saturday || (w == time.Sunday && d == 3) {
		commodities.SyncDailyCommoditiesData()
		commodities.SyncDailyEconomicData()
		count += 5
	}

	symbol, err := base.GetCurrentStock()
	if err != nil || len(symbol) == 0 {
		log.Println(err)
		return
	}

	// This will consume 2 tokens
	stock.SyncStockSymbolData(symbol)
	// This will consume 1 tokens
	news.SyncNewsSentimentData(symbol)
	count += 3

	for i := count; i < totalCount; i++ {
		stockSeries.SyncStockSeriesData(symbol)
	}

	base.UpdateStockRunningData(symbol)
}

func CommoditiesData() {
	commodities.SyncDailyCommoditiesData()
	commodities.SyncMonthlyCommoditiesData()
}

func EconomicIndicatorsData() {
	commodities.SyncDailyEconomicData()
	commodities.SyncMonthlyEconomicData()
	commodities.SyncQuarterlyEconomicData()
	commodities.SyncAnnualEconomicData()
}
