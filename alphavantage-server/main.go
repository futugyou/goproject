package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyou/alphavantage"
)

func main() {
	apikey := os.Getenv("ALPHAVANTAGE_API_KEY")
	fmt.Println(apikey)
	s := alphavantage.NewTimeSeriesClient(apikey)
	p := alphavantage.TimeSeriesParameter{
		Function: "TIME_SERIES_WEEKLY",
		Symbol:   "IBM",
		Interval: "15min",
	}
	result, err := s.ReadTimeSeries(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v.DividendAmount, v.SplitCoefficient)
	}
}
