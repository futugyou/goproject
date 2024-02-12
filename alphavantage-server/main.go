package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyou/alphavantage-server/stock"
)

func main() {
	stock.SyncStockSymbolData()
}
