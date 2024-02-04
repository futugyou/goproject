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

	// TimeSeries(s, dic)
	// TimeSeriesAdjusted(s, dic)
	// GlobalQuote(s, dic)
	// SymbolSearch(s, dic)
	MarketStatus(s, dic)
}

func TimeSeries(s *alphavantage.TimeSeriesClient, dic map[string]string) {
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
}

func TimeSeriesAdjusted(s *alphavantage.TimeSeriesClient, dic map[string]string) {
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

func GlobalQuote(s *alphavantage.TimeSeriesClient, dic map[string]string) {
	pp := alphavantage.GlobalQuoteParameter{
		Symbol: "IBM",
	}

	result1, err := s.GlobalQuote(pp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result1 {
		fmt.Println(v.Symbol, v.Open, v.High, v.Low, v.Price, v.Volume, v.LatestDay, v.PreviousClose, v.Change, v.ChangePercent)
	}
}

func SymbolSearch(s *alphavantage.TimeSeriesClient, dic map[string]string) {
	pp := alphavantage.SymbolSearchParameter{
		Keywords: "IBM",
	}

	result1, err := s.SymbolSearch(pp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result1 {
		fmt.Println(v.Symbol, v.Currency, v.MarketClose, v.MarketOpen, v.MatchScore, v.Name, v.Region, v.Timezone, v.Type)
	}
}

func MarketStatus(s *alphavantage.TimeSeriesClient, dic map[string]string) {

	v, err := s.MarketStatus()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(v.Endpoint)
	for _, vv := range v.Markets {
		fmt.Println(vv.CurrentStatus, vv.LocalClose, vv.LocalOpen, vv.MarketType, vv.Notes, vv.Notes, vv.PrimaryExchanges, vv.Region)
	}
}
