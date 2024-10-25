package api

import (
	_ "github.com/joho/godotenv/autoload"

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

	ctx := r.Context()
	ticker := r.URL.Query().Get("ticker")
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	repository := news.NewNewsRepository(config)
	var datas []news.NewsEntity
	var err error

	if ticker == "" {
		page := core.Paging{
			Page:      1,
			Limit:     1000,
			SortField: "time_published",
			Direct:    "DESC",
		}
		datas, err = repository.Paging(ctx, page)
	} else {
		datas, err = repository.GetAllByFilter(ctx, []core.DataFilterItem{{
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
