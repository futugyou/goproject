package alphavantage

import (
	"fmt"
	"strconv"
	"time"
)

const _Alphavantage_Http_Scheme string = "https"
const _Alphavantage_Host string = "www.alphavantage.co"
const _Alphavantage_Path string = "query"
const _Alphavantage_Datatype string = "csv"

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

func parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
