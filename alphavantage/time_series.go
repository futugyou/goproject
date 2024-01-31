package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
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

type TimeSeriesParameter struct {
	Function string `json:"function"`
	Symbol   string `json:"symbol"`
	Interval string `json:"interval"`
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
		datatype:   "csv",
	}
}

func (t *TimeSeriesClient) createRequestUrl(p TimeSeriesParameter) string {
	endpoint := &url.URL{}
	endpoint.Scheme = "https"
	endpoint.Host = "www.alphavantage.co"
	endpoint.Path = "query"
	query := endpoint.Query()
	query.Set("function", p.Function)
	query.Set("symbol", p.Symbol)
	query.Set("interval", p.Interval)
	query.Set("apikey", t.apikey)
	query.Set("datatype", t.datatype)
	endpoint.RawQuery = query.Encode()

	return endpoint.String()
}

func (t *TimeSeriesClient) ReadTimeSeries(p TimeSeriesParameter) ([]*TimeSeries, error) {
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

	d, err := time.Parse("2006-01-02 15:04:05", s[timestamp])
	if err != nil {
		return nil, fmt.Errorf("error parsing timestamp %s", s[timestamp])
	}
	value.Time = d

	f, err := strconv.ParseFloat(s[open], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing open %s", s[open])
	}
	value.Open = f

	f, err = strconv.ParseFloat(s[high], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing high %s", s[high])
	}
	value.High = f

	f, err = strconv.ParseFloat(s[low], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing low %s", s[low])
	}
	value.Low = f

	f, err = strconv.ParseFloat(s[close], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing close %s", s[close])
	}
	value.Close = f

	f, err = strconv.ParseFloat(s[volume], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing volume %s", s[volume])
	}
	value.Volume = f

	return value, nil
}
