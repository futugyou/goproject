package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage-server/news"
)

func News(w http.ResponseWriter, r *http.Request) {
	if crosForVercel(w, r) {
		return
	}
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := news.NewNewsRepository(config)
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
