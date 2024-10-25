package api

import (
	"context"
	"fmt"
	"html"

	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"net/http"
	"os"

	"github.com/futugyou/alphavantage-server/balance"
	"github.com/futugyou/alphavantage-server/cash"
	"github.com/futugyou/alphavantage-server/core"
	"github.com/futugyou/alphavantage-server/earnings"
	"github.com/futugyou/alphavantage-server/expected"
	"github.com/futugyou/alphavantage-server/income"

	"github.com/futugyou/extensions"
)

func Fundamentals(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	ctx := r.Context()
	datatype := r.URL.Query().Get("type")
	config := core.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}
	switch datatype {
	case "balance":
		balanceData(ctx, config, w)
	case "cash":
		cashData(ctx, config, w)
	case "earnings":
		earningsData(ctx, config, w)
	case "expected":
		expectedData(ctx, config, w)
	case "income":
		incomeData(ctx, config, w)
	default:
		fmt.Fprintf(w, "datatype %q is not support", html.EscapeString(datatype))
		w.WriteHeader(500)
		return
	}

}

func balanceData(ctx context.Context, config core.DBConfig, w http.ResponseWriter) {
	repository := balance.NewBalanceRepository(config)
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

func cashData(ctx context.Context, config core.DBConfig, w http.ResponseWriter) {
	repository := cash.NewCashRepository(config)
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

func earningsData(ctx context.Context, config core.DBConfig, w http.ResponseWriter) {
	repository := earnings.NewEarningsRepository(config)
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

func expectedData(ctx context.Context, config core.DBConfig, w http.ResponseWriter) {
	repository := expected.NewExpectedRepository(config)
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

func incomeData(ctx context.Context, config core.DBConfig, w http.ResponseWriter) {
	repository := income.NewIncomeRepository(config)
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
