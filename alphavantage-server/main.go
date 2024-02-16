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
	// StockSeries("IBM")
}

func ProcessToRun() {
	w := time.Now().Weekday()
	if w == time.Saturday {
		// Commodities data
		// Economic Indicators data
	} else {
		index := int(time.Sunday)
		symbol := stock.StockList[index]
		StockSeries(symbol)
		News(symbol)
	}

}

func SymbolData() {
	stock.SyncStockSymbolData()
}

func StockSeries(symbol string) {
	stockSeries.SyncStockSeriesData(symbol)
}

func News(symbol string) {
	news.SyncNewsSentimentData(symbol  )
}
