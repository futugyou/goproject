package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"
	"os"

	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage-server/stock"

	"github.com/futugyou/extensions"
)

func Company(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	getAllCompany(w, r)
}

func getAllCompany(w http.ResponseWriter, r *http.Request) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	ctx := r.Context()
	repository := stock.NewStockRepository(config)
	datas, err := repository.GetAll(ctx)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, _ := json.Marshal(datas)
	w.Write(body)
	w.WriteHeader(200)
}
