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

func AddNewStock(symbol string) {
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
	err := repository.Update(context.Background(), data, []core.DataFilterItem{{Key: "symbol", Value: symbol}})
	if err != nil {
		log.Println(err)
		return
	}

	stock.SyncStockSymbolData(symbol)
}

func UpdateStockRunningData(symbol string) {
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
	r, _ := repository.GetOne(context.Background(), []core.DataFilterItem{{Key: "symbol", Value: symbol}})
	if r != nil {
		data.RunCount = r.RunCount + 1
	}

	err := repository.Update(context.Background(), data, []core.DataFilterItem{{Key: "symbol", Value: symbol}})
	if err != nil {
		log.Println(err)
	}

	log.Printf("update %s running data finish. \n", symbol)
}

func GetCurrentStock() (string, error) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}
	repository := NewBaseDataRepository(config)
	data, err := repository.GetAll(context.Background())
	if err != nil || len(data) == 0 {
		log.Println(err)
		return "", err
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].RunCount < data[j].RunCount
	})

	return data[0].Symbol, nil
}

func InitAllStock() (bool, []string, error) {
	result := make([]string, 0)
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := NewBaseDataRepository(config)
	data, err := repository.GetAll(context.Background())
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
		AddNewStock(symbol)
		result = append(result, symbol)
	}

	return true, result, nil
}
