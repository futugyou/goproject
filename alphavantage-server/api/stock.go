package api

import (
	"encoding/json"
	"net/http"

	stockSeries "github.com/futugyou/alphavantage-server/stock-series"
)

func Stock(w http.ResponseWriter, r *http.Request) {
	if crosForVercel(w, r) {
		return
	}

	datas, err := stockSeries.StockSeriesData()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, _ := json.Marshal(datas)
	w.Write(body)
	w.WriteHeader(200)
}
