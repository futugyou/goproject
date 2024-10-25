package main

import (
	"context"
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
}

func ProcessToRun() {
	ctx := context.Background()
	// only init data when first time
	init, stockList, err := InitBaseData(ctx)
	if init || err != nil {
		return
	}

	d := time.Now().Day()
	w := time.Now().Weekday()
	count := 0
	// The number of available tokens is 25, but leave 3 as preparation
	totalCount := 22

	// 1. Runs on the 2nd of every month
	// This will consume 17 tokens, so return.
	if d == 2 {
		commodities.SyncMonthlyCommoditiesData(ctx)
		commodities.SyncMonthlyEconomicData(ctx)
		commodities.SyncQuarterlyEconomicData(ctx)
		commodities.SyncAnnualEconomicData(ctx)
		return
	}

	// Runs monthly stock sync from 3 to 3+len(StockList)
	// This will consume 5 tokens, so go on.
	for i := 0; i < len(stockList); i++ {
		if d == i+3 {
			symbol := stockList[i]
			income.SyncIncomeStatementData(ctx, symbol)
			balance.SyncBalanceSheetData(ctx, symbol)
			cash.SyncCashSheetData(ctx, symbol)
			earnings.SyncEarningsData(ctx, symbol)
			expected.SyncExpectedData(ctx, symbol)
			count += 5
		}
	}

	// Runs every Saturday or (3 of month and Sunday)
	// This will consume 5 tokens, so go on.
	if w == time.Saturday || (w == time.Sunday && d == 3) {
		commodities.SyncDailyCommoditiesData(ctx)
		commodities.SyncDailyEconomicData(ctx)
		count += 5
	}

	symbol, err := base.GetCurrentStock(ctx)
	if err != nil || len(symbol) == 0 {
		log.Println(err)
		return
	}

	// This will consume 2 tokens
	stock.SyncStockSymbolData(ctx, symbol)
	// This will consume 1 tokens
	news.SyncNewsSentimentData(ctx, symbol)
	count += 3

	for i := count; i < totalCount; i++ {
		if stockSeries.SyncStockSeriesData(ctx, symbol) {
			break
		}
	}

	base.UpdateStockRunningData(ctx, symbol)
}

func CommoditiesData(ctx context.Context) {
	commodities.SyncDailyCommoditiesData(ctx)
	commodities.SyncMonthlyCommoditiesData(ctx)
}

func EconomicIndicatorsData(ctx context.Context) {
	commodities.SyncDailyEconomicData(ctx)
	commodities.SyncMonthlyEconomicData(ctx)
	commodities.SyncQuarterlyEconomicData(ctx)
	commodities.SyncAnnualEconomicData(ctx)
}

func InitBaseData(ctx context.Context) (bool, []string, error) {
	init, list, err := base.InitAllStock(ctx)
	if init && err == nil {
		commodities.CreateCommoditiesIndex(ctx)
	}
	return init, list, err
}
