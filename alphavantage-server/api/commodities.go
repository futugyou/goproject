package api

import (
	_ "github.com/joho/godotenv/autoload"

	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/futugyou/alphavantage-server/commodities"
	"github.com/futugyou/alphavantage-server/core"
)

func Commodities(w http.ResponseWriter, r *http.Request) {
	if crosForVercel(w, r) {
		return
	}

	datatype := r.URL.Query().Get("type")
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}
	repository := commodities.NewCommoditiesRepository(config)
	var datas []commodities.CommoditiesEntity
	var err error
	if len(datatype) > 0 && datatype != "ALL" {
		datas, err = repository.GetCommoditiesByType(context.Background(), datatype)
	} else {
		datas, err = repository.GetAll(context.Background())
	}
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, _ := json.Marshal(datas)
	w.Write(body)
	w.WriteHeader(200)

}
