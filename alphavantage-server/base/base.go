package base

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/futugyou/alphavantage-server/core"
)

func AddNewStock(symbol string) {
	log.Printf("begin to add %s to db. \n", symbol)
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	data := []BaseDataEntity{{
		Id:       symbol,
		Symbol:   symbol,
		RunCount: 0,
		RunDate:  time.Now(),
	}}
	repository := NewBaseDataRepository(config)
	r, err := repository.InsertMany(context.Background(), data, BalanceFilter)
	if err != nil {
		log.Println(err)
		return
	}

	r.String()
	log.Println("balance data sync finish")
}

func BalanceFilter(e BaseDataEntity) []core.DataFilterItem {
	return []core.DataFilterItem{{Key: "symbol", Value: e.Symbol}}
}
