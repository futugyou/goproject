package stock

import (
	"context"
	"os"

	"log"

	"github.com/futugyou/alphavantage"
	"github.com/futugyou/alphavantage-server/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SyncStockSymbolData() {
	// 1. get data from api
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	client := alphavantage.NewTimeSeriesClient(apikey)
	symbolList := make([]*alphavantage.SymbolSearch, 0)
	for _, key := range StockList {
		pp := alphavantage.SymbolSearchParameter{
			Keywords: key,
		}

		result, err := client.SymbolSearch(pp)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(result) > 0 {
			symbolList = append(symbolList, result[0])
		}
	}

	if len(symbolList) == 0 {
		return
	}

	// 2. create insert data
	inputList := make([]StockEntity, len(symbolList))
	for i := 0; i < len(symbolList); i++ {
		inputList[i] = StockEntity{
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

	// 3. insert data
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewStockRepository(config)
	repository.InsertMany(context.Background(), inputList, StockSymbolFilter)
}

func StockSymbolFilter(e StockEntity) primitive.D {
	return bson.D{{Key: "symbol", Value: e.Symbol}}
}
