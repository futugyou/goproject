package main

import "fmt"

func main() {
	h := NewHttpClient()
	result, _ := h.Get("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=IBM&interval=5min&apikey=demo&datatype=csv")
	switch result := result.(type) {
	case []*TimeSeries:
		for _, v := range result {
			fmt.Println(v.Close)
		}
	default:
		fmt.Println("something error")
	}
}
