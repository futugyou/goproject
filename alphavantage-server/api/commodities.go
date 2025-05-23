package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"
	"os"

	"github.com/futugyou/alphavantage-server/commodities"
	"github.com/futugyou/alphavantage-server/core"

	"github.com/futugyou/extensions"
)

func Commodities(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	ctx := r.Context()
	datatype := r.URL.Query().Get("type")
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}
	repository := commodities.NewCommoditiesRepository(config)
	var datas []commodities.CommoditiesEntity
	var err error
	if len(datatype) > 0 && datatype != "ALL" {
		datas, err = repository.GetCommoditiesByType(ctx, datatype)
	} else {
		datas, err = repository.GetAll(ctx)
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
