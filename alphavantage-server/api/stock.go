package api

import (
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"

	stockSeries "github.com/futugyou/alphavantage-server/stock-series"

	"github.com/futugyou/extensions"
)

func Stock(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	symbol := r.URL.Query().Get("symbol")
	year := r.URL.Query().Get("year")
	datas, err := stockSeries.StockSeriesData(symbol, year)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	body, _ := json.Marshal(datas)
	w.Write(body)
	w.WriteHeader(200)
}
