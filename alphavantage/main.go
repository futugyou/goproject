package main

import (
	_ "github.com/joho/godotenv/autoload"

	"fmt"
)

func main() {
	h := NewHttpClient()
	s := NewTimeSeriesClient(h)

	result, err := s.ReadTimeSeries()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range result {
		fmt.Println(v.Close)
	}
}
