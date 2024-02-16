package main

import (
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyou/alphavantage-server/news"
	"github.com/futugyou/alphavantage-server/stock"
	stockSeries "github.com/futugyou/alphavantage-server/stock-series"
)

func main() {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		ProcessToRun()
		return
	}
	SymbolData("IBM")
}

func ProcessToRun() {
	w := time.Now().Weekday()
	if w == time.Saturday {
		// Commodities data
		// Economic Indicators data
	} else {
		index := int(time.Sunday)
		symbol := stock.StockList[index]
		SymbolData(symbol)
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
