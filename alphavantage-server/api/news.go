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
	filter := []core.DataFilterItem{}
	if len(ticker) > 0 {
		filter = append(filter, core.DataFilterItem{
			Key:   "ticker_sentiment",
			Value: map[string]any{"$elemMatch": map[string]any{"$eq": ticker}},
		})
	}
	page := core.Paging{
		Page:      1,
		Limit:     1000,
		SortField: "time_published",
		Direct:    "DESC",
	}
	datas, err := repository.GetWithFilterAndPaging(ctx, filter, &page)

	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, _ := json.Marshal(datas)
	w.Write(body)
	w.WriteHeader(200)
}
