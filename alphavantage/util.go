package main

import (
	"fmt"
	"time"
)

const Alphavantage_Http_Scheme string = "https"
const Alphavantage_Host string = "www.alphavantage.co"
const Alphavantage_Path string = "query"
const Alphavantage_Datatype string = "csv"

var time_format = []string{"2006-01-02 15:04:05", "2006-01-02"}

func parseTime(value string) (time.Time, error) {
	for _, f := range time_format {
		d, err := time.Parse(f, value)
		if err == nil {
			return d, nil
		}
	}

	return time.Time{}, fmt.Errorf("time parse error, raw data is %s", value)
}
