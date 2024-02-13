package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyou/alphavantage-server/stock"
	stockSeries "github.com/futugyou/alphavantage-server/stock-series"
)

func main() {
	// SymbolData()
	StockSeries()
}

func SymbolData() {
	stock.SyncStockSymbolData()
}

func StockSeries() {
	stockSeries.SyncStockSeriesData()
}
