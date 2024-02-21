package api

import (
	"context"
	"os"

	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage-server/stock"
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"
)

func crosForVercel(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, x-requested-with, account-id")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin, Token, Content-Length, Access-Control-Allow-Headers, Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}

	return false
}

func Company(w http.ResponseWriter, r *http.Request) {
	if crosForVercel(w, r) {
		return
	}

	getAllCompany(w, r)
}

func getAllCompany(w http.ResponseWriter, r *http.Request) {
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := stock.NewStockRepository(config)
	datas, err := repository.GetAll(context.Background())
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, _ := json.Marshal(datas)
	w.Write(body)
	w.WriteHeader(200)
}
