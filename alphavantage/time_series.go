package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"slices"
	"time"
)

type TimeSeries struct {
	Symbol string    `json:"symbol"`
	Time   time.Time `json:"time"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume float64   `json:"volume"`
}

// const (
// 	TIME_SERIES_INTRADAY = iota
// 	TIME_SERIES_DAILY
// 	TIME_SERIES_DAILY_ADJUSTED
// 	TIME_SERIES_WEEKLY
// 	TIME_SERIES_WEEKLY_ADJUSTED
// 	TIME_SERIES_MONTHLY
// 	TIME_SERIES_MONTHLY_ADJUSTED
// 	GLOBAL_QUOTE
// 	SYMBOL_SEARCH
// 	MARKET_STATUS
// )

type timeSeriesFunctionType struct {
	name string
}

// const TIME_SERIES_INTRADAY FunctionType = "TIME_SERIES_INTRADAY"
// const TIME_SERIES_DAILY FunctionType = "TIME_SERIES_DAILY"
// const TIME_SERIES_DAILY_ADJUSTED FunctionType = "TIME_SERIES_DAILY_ADJUSTED"
// const TIME_SERIES_WEEKLY FunctionType = "TIME_SERIES_WEEKLY"
// const TIME_SERIES_WEEKLY_ADJUSTED FunctionType = "TIME_SERIES_WEEKLY_ADJUSTED"
// const TIME_SERIES_MONTHLY FunctionType = "TIME_SERIES_MONTHLY"
// const TIME_SERIES_MONTHLY_ADJUSTED FunctionType = "TIME_SERIES_MONTHLY_ADJUSTED"
// const GLOBAL_QUOTE FunctionType = "GLOBAL_QUOTE"
// const SYMBOL_SEARCH FunctionType = "SYMBOL_SEARCH"
// const MARKET_STATUS FunctionType = "MARKET_STATUS"

type TimeSeriesParameter struct {
	Function   string            `json:"function"`
	Symbol     string            `json:"symbol"`
	Interval   string            `json:"interval"`
	Dictionary map[string]string `json:"dictionary"`
}

type TimeSeriesClient struct {
	httpClient *httpClient
	apikey     string
	datatype   string
}

func NewTimeSeriesClient() *TimeSeriesClient {
	return &TimeSeriesClient{
		httpClient: NewHttpClient(),
		apikey:     os.Getenv("ALPHAVANTAGE_API_KEY"),
		datatype:   Alphavantage_Datatype,
	}
}

func (t *TimeSeriesClient) createRequestUrl(p TimeSeriesParameter) string {
	endpoint := &url.URL{}
	endpoint.Scheme = Alphavantage_Http_Scheme
	endpoint.Host = Alphavantage_Host
	endpoint.Path = Alphavantage_Path
	query := endpoint.Query()
	query.Set("function", p.Function)
	query.Set("symbol", p.Symbol)
	query.Set("interval", p.Interval)
	query.Set("apikey", t.apikey)
	query.Set("datatype", t.datatype)
	for k, v := range p.Dictionary {
		query.Set(k, v)
	}
	endpoint.RawQuery = query.Encode()

	return endpoint.String()
}

func (t *TimeSeriesClient) ReadTimeSeries(p TimeSeriesParameter) ([]*TimeSeries, error) {
	_, err := t.checkTimeSeriesParameter(p.Function)
	if err != nil {
		return nil, err
	}

	path := t.createRequestUrl(p)
	r, err := t.httpClient.get(path)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	reader := csv.NewReader(r)
	reader.ReuseRecord = true
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	result := make([]*TimeSeries, 0)

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		value, err := t.readTimeSeriesItem(record)
		if err != nil {
			return nil, err
		}

		value.Symbol = p.Symbol
		result = append(result, value)
	}

	return result, nil

}

func (t *TimeSeriesClient) readTimeSeriesItem(s []string) (*TimeSeries, error) {
	const (
		timestamp = iota
		open
		high
		low
		close
		volume
	)

	value := &TimeSeries{}

	d, err := parseTime(s[timestamp])
	if err != nil {
		return nil, fmt.Errorf("error parsing timestamp %s", s[timestamp])
	}
	value.Time = d

	f, err := parseFloat(s[open])
	if err != nil {
		return nil, fmt.Errorf("error parsing open %s", s[open])
	}
	value.Open = f

	f, err = parseFloat(s[high])
	if err != nil {
		return nil, fmt.Errorf("error parsing high %s", s[high])
	}
	value.High = f

	f, err = parseFloat(s[low])
	if err != nil {
		return nil, fmt.Errorf("error parsing low %s", s[low])
	}
	value.Low = f

	f, err = parseFloat(s[close])
	if err != nil {
		return nil, fmt.Errorf("error parsing close %s", s[close])
	}
	value.Close = f

	f, err = parseFloat(s[volume])
	if err != nil {
		return nil, fmt.Errorf("error parsing volume %s", s[volume])
	}
	value.Volume = f

	return value, nil
}

func (t *TimeSeriesClient) checkTimeSeriesParameter(function string) (*timeSeriesFunctionType, error) {
	if have := slices.Contains(functionList, function); !have {
		return nil, fmt.Errorf("invalid function name %s", function)
	}

	result := &timeSeriesFunctionType{
		name: function,
	}
	return result, nil
}

var functionList = []string{"TIME_SERIES_INTRADAY", "TIME_SERIES_DAILY", "TIME_SERIES_DAILY_ADJUSTED",
	"TIME_SERIES_WEEKLY", "TIME_SERIES_WEEKLY_ADJUSTED", "TIME_SERIES_MONTHLY", "TIME_SERIES_MONTHLY_ADJUSTED",
	"GLOBAL_QUOTE", "SYMBOL_SEARCH", "MARKET_STATUS",
}
