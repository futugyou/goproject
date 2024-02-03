package alphavantage

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const _TIME_SERIES_INTRADAY string = "TIME_SERIES_INTRADAY"
const _TIME_SERIES_DAILY string = "TIME_SERIES_DAILY"
const _TIME_SERIES_DAILY_ADJUSTED string = "TIME_SERIES_DAILY_ADJUSTED"
const _TIME_SERIES_WEEKLY string = "TIME_SERIES_WEEKLY"
const _TIME_SERIES_WEEKLY_ADJUSTED string = "TIME_SERIES_WEEKLY_ADJUSTED"
const _TIME_SERIES_MONTHLY string = "TIME_SERIES_MONTHLY"
const _TIME_SERIES_MONTHLY_ADJUSTED string = "TIME_SERIES_MONTHLY_ADJUSTED"
const _GLOBAL_QUOTE string = "GLOBAL_QUOTE"
const _SYMBOL_SEARCH string = "SYMBOL_SEARCH"
const _MARKET_STATUS string = "MARKET_STATUS"

const _1min string = "1min"
const _5min string = "5min"
const _15min string = "15min"
const _30min string = "30min"
const _60min string = "60min"

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
	if value[len(value)-1] == '%' {
		v, err := strconv.ParseFloat(strings.Trim(value, "%"), 64)
		if err != nil {
			return 0, err
		}
		return v / 100, nil
	}
	return strconv.ParseFloat(value, 64)
}
