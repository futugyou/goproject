package api

import (
	_ "github.com/joho/godotenv/autoload"

	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage-server/news"

	"github.com/futugyou/extensions"
)

func News(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	ticker := r.URL.Query().Get("ticker")
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := news.NewNewsRepository(config)
	var datas []news.NewsEntity
	var err error

	if ticker == "" {
		datas, err = repository.GetAll(context.Background())
	} else {
		datas, err = repository.GetAllByFilter(context.Background(), []core.DataFilterItem{{
			Key:   "ticker_sentiment",
			Value: map[string]interface{}{"$elemMatch": map[string]interface{}{"$eq": ticker}},
		}})
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
