package base

import (
	"context"
	"log"
	"os"
	"sort"
	"time"

	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage-server/stock"
)

func AddNewStock(ctx context.Context, symbol string) {
	log.Printf("begin to add %s to db. \n", symbol)
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	data := BaseDataEntity{
		Id:       symbol,
		Symbol:   symbol,
		RunCount: 0,
		RunDate:  time.Now(),
	}

	repository := NewBaseDataRepository(config)
	err := repository.Update(ctx, data, []core.DataFilterItem{{Key: "symbol", Value: symbol}})
	if err != nil {
		log.Println(err)
		return
	}

	stock.SyncStockSymbolData(ctx, symbol)
}

func UpdateStockRunningData(ctx context.Context, symbol string) {
	log.Printf("begin to update %s running data. \n", symbol)
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	data := BaseDataEntity{
		Id:       symbol,
		Symbol:   symbol,
		RunCount: 0,
		RunDate:  time.Now(),
	}

	repository := NewBaseDataRepository(config)
	r, _ := repository.GetOne(ctx, []core.DataFilterItem{{Key: "symbol", Value: symbol}})
	if r != nil {
		data.RunCount = r.RunCount + 1
	}

	err := repository.Update(ctx, data, []core.DataFilterItem{{Key: "symbol", Value: symbol}})
	if err != nil {
		log.Println(err)
	}

	log.Printf("update %s running data finish. \n", symbol)
}

func GetCurrentStock(ctx context.Context) (string, error) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}
	repository := NewBaseDataRepository(config)
	data, err := repository.GetAll(ctx)
	if err != nil || len(data) == 0 {
		log.Println(err)
		return "", err
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].RunCount < data[j].RunCount
	})

	return data[0].Symbol, nil
}

func InitAllStock(ctx context.Context) (bool, []string, error) {
	result := make([]string, 0)
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewBaseDataRepository(config)
	data, err := repository.GetAll(ctx)
	if err != nil {
		log.Println(err)
		return false, result, err
	}

	if len(data) > 0 {
		for i := 0; i < len(data); i++ {
			result = append(result, data[i].Symbol)
		}
		return false, result, nil
	}

	for i := 0; i < len(stock.StockList); i++ {
		symbol := stock.StockList[i]
		AddNewStock(ctx, symbol)
		result = append(result, symbol)
	}

	return true, result, nil
}
