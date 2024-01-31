package main

import (
	_ "github.com/joho/godotenv/autoload"

	"fmt"
)

func main() {
	s := NewTimeSeriesClient()
	p := TimeSeriesParameter{
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
