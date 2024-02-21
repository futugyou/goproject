package api

import (
	"encoding/json"
	"net/http"

	"github.com/futugyou/alphavantage-server/news"
)

func News(w http.ResponseWriter, r *http.Request) {
	if crosForVercel(w, r) {
		return
	}

	datas, err := news.NewsData()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, _ := json.Marshal(datas)
	w.Write(body)
	w.WriteHeader(200)
}
