package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyou/alphavantage"
)

func main() {
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")

	dic := make(map[string]string)
	dic["month"] = "2024-01"

	s := alphavantage.NewTimeSeriesClient(apikey)
	p := alphavantage.TimeSeriesIntradayParameter{
		Symbol:     "IBM",
		Interval:   "15min",
		Dictionary: dic,
	}
	result, err := s.TimeSeriesIntraday(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v.Symbol, v.Time, v.Open, v.High, v.Low, v.Close, v.Volume)
	}

	pp := alphavantage.TimeSeriesMonthlyAdjustedParameter{
		Symbol:     "IBM",
		Dictionary: dic,
	}

	result1, err := s.TimeSeriesMonthlyAdjusted(pp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result1 {
		fmt.Println(v.Symbol, v.Time, v.Open, v.High, v.Low, v.Close, v.Volume, v.AdjustedClose, v.DividendAmount, v.SplitCoefficient)
	}
}
