package alphavantage

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/futugyou/alphavantage/enums"
)

// timestamp,open,high,low,close,volume
type TimeSeries struct {
	Symbol string    `json:"symbol"`
	Time   time.Time `json:"time"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume float64   `json:"volume"`
}

// timestamp,open,high,low,close,adjusted_close,volume,dividend_amount,split_coefficient
type TimeSeriesAdjusted struct {
	Symbol           string    `json:"symbol"`
	Time             time.Time `json:"time"`
	Open             float64   `json:"open"`
	High             float64   `json:"high"`
	Low              float64   `json:"low"`
	Close            float64   `json:"close"`
	Volume           float64   `json:"volume"`
	AdjustedClose    float64   `json:"adjusted_close"`
	DividendAmount   float64   `json:"dividend_amount"`
	SplitCoefficient float64   `json:"split_coefficient"`
}

type timeSeriesParameter struct {
	Function   string             `json:"function"`
	Symbol     string             `json:"symbol"`
	Interval   enums.TimeInterval `json:"interval"`
	Dictionary map[string]string  `json:"dictionary"`
}

// parameter for TIME_SERIES_INTRADAY API
type TimeSeriesIntradayParameter struct {
	// The name of the equity of your choice. For example: symbol=IBM
	Symbol string `json:"symbol"`
	// Time interval between two consecutive data points in the time series. The following values are supported: 1min, 5min, 15min, 30min, 60min
	Interval enums.TimeInterval `json:"interval"`
	// other option parameter, see https://www.alphavantage.co/documentation/#intraday
	Dictionary map[string]string `json:"dictionary"`
}

func (t TimeSeriesIntradayParameter) Validation() error {
	if len(strings.TrimSpace(t.Symbol)) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}
	return nil
}

// parameter for TIME_SERIES_DAILY API
type TimeSeriesDailyParameter struct {
	// The name of the equity of your choice. For example: symbol=IBM
	Symbol string `json:"symbol"`
	// other option parameter, see https://www.alphavantage.co/documentation/#daily
	Dictionary map[string]string `json:"dictionary"`
}

func (t TimeSeriesDailyParameter) Validation() error {
	if len(strings.TrimSpace(t.Symbol)) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}

	return nil
}

// parameter for TIME_SERIES_WEEKLY API
type TimeSeriesWeeklyParameter struct {
	// The name of the equity of your choice. For example: symbol=IBM
	Symbol string `json:"symbol"`
	// other option parameter, see https://www.alphavantage.co/documentation/#weekly
	Dictionary map[string]string `json:"dictionary"`
}

func (t TimeSeriesWeeklyParameter) Validation() error {
	if len(strings.TrimSpace(t.Symbol)) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}

	return nil
}

// parameter for TIME_SERIES_MONTHLY API
type TimeSeriesMonthlyParameter struct {
	// The name of the equity of your choice. For example: symbol=IBM
	Symbol string `json:"symbol"`
	// other option parameter, see https://www.alphavantage.co/documentation/#monthly
	Dictionary map[string]string `json:"dictionary"`
}

func (t TimeSeriesMonthlyParameter) Validation() error {
	if len(strings.TrimSpace(t.Symbol)) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}

	return nil
}

type TimeSeriesClient struct {
	httpClient *httpClient
	apikey     string
	datatype   string
}

// api client for get Time Series Stock Data
// This suite of APIs provide global equity data in 4 different temporal resolutions:
// (1) daily, (2) weekly, (3) monthly, and (4) intraday, with 20+ years of historical depth.
// A lightweight ticker quote endpoint and several utility functions
// such as ticker search and market open/closure status are also included for your convenience.
func NewTimeSeriesClient(apikey string) *TimeSeriesClient {
	return &TimeSeriesClient{
		httpClient: newHttpClient(),
		apikey:     apikey,
		datatype:   _Alphavantage_Datatype,
	}
}

// This API returns current and 20+ years of historical intraday OHLCV time series of the equity specified,
// covering extended trading hours where applicable (e.g., 4:00am to 8:00pm Eastern Time for the US market).
// You can query both raw (as-traded) and split/dividend-adjusted intraday data from this endpoint.
// The OHLCV data is sometimes called "candles" in finance literature.
func (t *TimeSeriesClient) TimeSeriesIntraday(p TimeSeriesIntradayParameter) ([]*TimeSeries, error) {
	err := p.Validation()
	if err != nil {
		return nil, err
	}

	innnerParameter := timeSeriesParameter{
		Function:   _TIME_SERIES_INTRADAY,
		Symbol:     p.Symbol,
		Interval:   p.Interval,
		Dictionary: p.Dictionary,
	}

	return t.readTimeSeries(innnerParameter)
}

// This API returns raw (as-traded) daily time series (date, daily open, daily high, daily low, daily close, daily volume) of the global equity specified,
// covering 20+ years of historical data. The OHLCV data is sometimes called "candles" in finance literature.
// If you are also interested in split/dividend-adjusted data, please use the Daily Adjusted API,
// which covers adjusted close values and historical split and dividend events.
func (t *TimeSeriesClient) TimeSeriesDaily(p TimeSeriesDailyParameter) ([]*TimeSeries, error) {
	err := p.Validation()
	if err != nil {
		return nil, err
	}

	innnerParameter := timeSeriesParameter{
		Function:   _TIME_SERIES_DAILY,
		Symbol:     p.Symbol,
		Dictionary: p.Dictionary,
	}

	return t.readTimeSeries(innnerParameter)
}

// This API returns weekly time series
// (last trading day of each week, weekly open, weekly high, weekly low, weekly close, weekly volume)
// of the global equity specified, covering 20+ years of historical data.
func (t *TimeSeriesClient) TimeSeriesWeekly(p TimeSeriesWeeklyParameter) ([]*TimeSeries, error) {
	err := p.Validation()
	if err != nil {
		return nil, err
	}

	innnerParameter := timeSeriesParameter{
		Function:   _TIME_SERIES_WEEKLY,
		Symbol:     p.Symbol,
		Dictionary: p.Dictionary,
	}

	return t.readTimeSeries(innnerParameter)
}

// This API returns monthly time series
// (last trading day of each month, monthly open, monthly high, monthly low, monthly close, monthly volume)
// of the global equity specified, covering 20+ years of historical data.
func (t *TimeSeriesClient) TimeSeriesMonthly(p TimeSeriesMonthlyParameter) ([]*TimeSeries, error) {
	err := p.Validation()
	if err != nil {
		return nil, err
	}

	innnerParameter := timeSeriesParameter{
		Function:   _TIME_SERIES_MONTHLY,
		Symbol:     p.Symbol,
		Dictionary: p.Dictionary,
	}

	return t.readTimeSeries(innnerParameter)
}

func (t *TimeSeriesClient) readTimeSeries(p timeSeriesParameter) ([]*TimeSeries, error) {
	path := t.createRequestUrl(p)
	csvData, err := t.httpClient.getCsv(path)
	if err != nil {
		return nil, err
	}

	result := make([]*TimeSeries, 0)

	for i := 0; i < len(csvData); i++ {
		value, err := t.readTimeSeriesItem(csvData[i])
		if err != nil {
			return nil, err
		}

		value.Symbol = p.Symbol
		result = append(result, value)
	}

	return result, nil
}

func (t *TimeSeriesClient) createRequestUrl(p timeSeriesParameter) string {
	endpoint := &url.URL{}
	endpoint.Scheme = _Alphavantage_Http_Scheme
	endpoint.Host = _Alphavantage_Host
	endpoint.Path = _Alphavantage_Path
	query := endpoint.Query()
	query.Set("function", p.Function)
	query.Set("symbol", p.Symbol)
	query.Set("interval", p.Interval.String())
	query.Set("apikey", t.apikey)
	query.Set("datatype", t.datatype)
	for k, v := range p.Dictionary {
		query.Set(k, v)
	}
	endpoint.RawQuery = query.Encode()

	return endpoint.String()
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
