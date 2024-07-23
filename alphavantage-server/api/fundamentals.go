package api

import (
	"fmt"
	"html"

	_ "github.com/joho/godotenv/autoload"

	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/futugyou/alphavantage-server/balance"
	"github.com/futugyou/alphavantage-server/cash"
	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage-server/earnings"
	"github.com/futugyou/alphavantage-server/expected"
	"github.com/futugyou/alphavantage-server/income"
)

func Fundamentals(w http.ResponseWriter, r *http.Request) {
	if crosForVercel(w, r) {
		return
	}

	datatype := r.URL.Query().Get("type")
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}
	switch datatype {
	case "balance":
		balanceData(config, w)
	case "cash":
		cashData(config, w)
	case "earnings":
		earningsData(config, w)
	case "expected":
		expectedData(config, w)
	case "income":
		incomeData(config, w)
	default:
		fmt.Fprintf(w, "datatype %q is not support", html.EscapeString(datatype))
		w.WriteHeader(500)
		return
	}

}

func balanceData(config core.DBConfig, w http.ResponseWriter) {
	repository := balance.NewBalanceRepository(config)
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

func cashData(config core.DBConfig, w http.ResponseWriter) {
	repository := cash.NewCashRepository(config)
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

func earningsData(config core.DBConfig, w http.ResponseWriter) {
	repository := earnings.NewEarningsRepository(config)
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

func expectedData(config core.DBConfig, w http.ResponseWriter) {
	repository := expected.NewExpectedRepository(config)
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

func incomeData(config core.DBConfig, w http.ResponseWriter) {
	repository := income.NewIncomeRepository(config)
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
