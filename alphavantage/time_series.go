package alphavantage

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"
)

var timeSeriesDataFunctionList = []string{_TIME_SERIES_INTRADAY, _TIME_SERIES_DAILY,
	_TIME_SERIES_WEEKLY, _TIME_SERIES_MONTHLY,
}

var timeSeriesDataaAjustedFunctionList = []string{_TIME_SERIES_DAILY_ADJUSTED,
	_TIME_SERIES_WEEKLY_ADJUSTED, _TIME_SERIES_MONTHLY_ADJUSTED,
}

var timeSeriesDataIntervalList = []string{_1min, _5min, _15min, _30min, _60min}

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

// symbol,open,high,low,price,volume,latestDay,previousClose,change,changePercent
type GlobalEndpoint struct {
	Symbol        string    `json:"symbol"`
	Open          float64   `json:"open"`
	High          float64   `json:"high"`
	Low           float64   `json:"low"`
	Price         float64   `json:"price"`
	Volume        float64   `json:"volume"`
	LatestDay     time.Time `json:"latestDay"`
	PreviousClose float64   `json:"previous_close"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"change_percent"`
}

// symbol,name,type,region,marketOpen,marketClose,timezone,currency,matchScore
type SymbolSearch struct {
	Symbol      string        `json:"symbol"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Region      string        `json:"region"`
	MarketOpen  time.Duration `json:"marketOpen"`
	MarketClose time.Duration `json:"marketClose"`
	Timezone    string        `json:"timezone"`
	Currency    string        `json:"currency"`
	MatchScore  float64       `json:"matchScore"`
}

type MarketStatus struct {
	Endpoint string   `json:"endpoint"`
	Markets  []Market `json:"markets"`
}

type Market struct {
	MarketType       string `json:"market_type"`
	Region           string `json:"region"`
	PrimaryExchanges string `json:"primary_exchanges"`
	LocalOpen        string `json:"local_open"`
	LocalClose       string `json:"local_close"`
	CurrentStatus    string `json:"current_status"`
	Notes            string `json:"notes"`
}

type timeSeriesFunctionType struct {
	name string
}

type timeSeriesParameter struct {
	Function   string            `json:"function"`
	Symbol     string            `json:"symbol"`
	Interval   string            `json:"interval"`
	Dictionary map[string]string `json:"dictionary"`
}

// parameter for TIME_SERIES_INTRADAY API
type TimeSeriesIntradayParameter struct {
	// The name of the equity of your choice. For example: symbol=IBM
	Symbol string `json:"symbol"`
	// Time interval between two consecutive data points in the time series. The following values are supported: 1min, 5min, 15min, 30min, 60min
	Interval string `json:"interval"`
	// other option parameter, see https://www.alphavantage.co/documentation/#intraday
	Dictionary map[string]string `json:"dictionary"`
}

func (t TimeSeriesIntradayParameter) Validation() error {
	if len(strings.Trim(t.Symbol, " ")) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}

	if have := slices.Contains(timeSeriesDataIntervalList, t.Interval); !have {
		return fmt.Errorf("invalid interval name %s, allowed interval are  %s", t.Interval, strings.Join(timeSeriesDataIntervalList, ","))
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
	if len(strings.Trim(t.Symbol, " ")) == 0 {
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
	if len(strings.Trim(t.Symbol, " ")) == 0 {
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
	if len(strings.Trim(t.Symbol, " ")) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}

	return nil
}

// parameter for TIME_SERIES_DAILY_ADJUSTED API
type TimeSeriesDailyAdjustedParameter struct {
	// The name of the equity of your choice. For example: symbol=IBM
	Symbol string `json:"symbol"`
	// other option parameter, see https://www.alphavantage.co/documentation/#dailyadj
	Dictionary map[string]string `json:"dictionary"`
}

func (t TimeSeriesDailyAdjustedParameter) Validation() error {
	if len(strings.Trim(t.Symbol, " ")) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}

	return nil
}

// parameter for TIME_SERIES_WEEKLY_ADJUSTED API
type TimeSeriesWeeklyAdjustedParameter struct {
	// The name of the equity of your choice. For example: symbol=IBM
	Symbol string `json:"symbol"`
	// other option parameter, see https://www.alphavantage.co/documentation/#weeklyadj
	Dictionary map[string]string `json:"dictionary"`
}

func (t TimeSeriesWeeklyAdjustedParameter) Validation() error {
	if len(strings.Trim(t.Symbol, " ")) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}

	return nil
}

// parameter for TIME_SERIES_MONTHLY_ADJUSTED API
type TimeSeriesMonthlyAdjustedParameter struct {
	// The name of the equity of your choice. For example: symbol=IBM
	Symbol string `json:"symbol"`
	// other option parameter, see https://www.alphavantage.co/documentation/#monthlyadj
	Dictionary map[string]string `json:"dictionary"`
}

func (t TimeSeriesMonthlyAdjustedParameter) Validation() error {
	if len(strings.Trim(t.Symbol, " ")) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}

	return nil
}

type timeSeriesClient struct {
	httpClient *httpClient
	apikey     string
	datatype   string
}

// api client for get Time Series Stock Data
// This suite of APIs provide global equity data in 4 different temporal resolutions:
// (1) daily, (2) weekly, (3) monthly, and (4) intraday, with 20+ years of historical depth.
// A lightweight ticker quote endpoint and several utility functions
// such as ticker search and market open/closure status are also included for your convenience.
func NewTimeSeriesClient(apikey string) *timeSeriesClient {
	return &timeSeriesClient{
		httpClient: newHttpClient(),
		apikey:     apikey,
		datatype:   _Alphavantage_Datatype,
	}
}

// This API returns current and 20+ years of historical intraday OHLCV time series of the equity specified,
// covering extended trading hours where applicable (e.g., 4:00am to 8:00pm Eastern Time for the US market).
// You can query both raw (as-traded) and split/dividend-adjusted intraday data from this endpoint.
// The OHLCV data is sometimes called "candles" in finance literature.
func (t *timeSeriesClient) TimeSeriesIntraday(p TimeSeriesIntradayParameter) ([]*TimeSeries, error) {
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
func (t *timeSeriesClient) TimeSeriesDaily(p TimeSeriesDailyParameter) ([]*TimeSeries, error) {
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
func (t *timeSeriesClient) TimeSeriesWeekly(p TimeSeriesWeeklyParameter) ([]*TimeSeries, error) {
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
func (t *timeSeriesClient) TimeSeriesMonthly(p TimeSeriesMonthlyParameter) ([]*TimeSeries, error) {
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

func (t *timeSeriesClient) readTimeSeries(p timeSeriesParameter) ([]*TimeSeries, error) {
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

func (t *timeSeriesClient) TimeSeriesDailyAdjusted(p TimeSeriesDailyAdjustedParameter) ([]*TimeSeriesAdjusted, error) {
	err := p.Validation()
	if err != nil {
		return nil, err
	}

	innnerParameter := timeSeriesParameter{
		Function:   _TIME_SERIES_DAILY_ADJUSTED,
		Symbol:     p.Symbol,
		Dictionary: p.Dictionary,
	}
	return t.readTimeSeriesAdjusted(innnerParameter)
}

func (t *timeSeriesClient) TimeSeriesWeeklyAdjusted(p TimeSeriesWeeklyAdjustedParameter) ([]*TimeSeriesAdjusted, error) {
	err := p.Validation()
	if err != nil {
		return nil, err
	}

	innnerParameter := timeSeriesParameter{
		Function:   _TIME_SERIES_WEEKLY_ADJUSTED,
		Symbol:     p.Symbol,
		Dictionary: p.Dictionary,
	}
	return t.readTimeSeriesAdjusted(innnerParameter)
}

func (t *timeSeriesClient) TimeSeriesMonthlyAdjusted(p TimeSeriesMonthlyAdjustedParameter) ([]*TimeSeriesAdjusted, error) {
	err := p.Validation()
	if err != nil {
		return nil, err
	}

	innnerParameter := timeSeriesParameter{
		Function:   _TIME_SERIES_MONTHLY_ADJUSTED,
		Symbol:     p.Symbol,
		Dictionary: p.Dictionary,
	}
	return t.readTimeSeriesAdjusted(innnerParameter)
}

func (t *timeSeriesClient) readTimeSeriesAdjusted(p timeSeriesParameter) ([]*TimeSeriesAdjusted, error) {
	path := t.createRequestUrl(p)
	csvData, err := t.httpClient.getCsv(path)
	if err != nil {
		return nil, err
	}

	result := make([]*TimeSeriesAdjusted, 0)

	for i := 0; i < len(csvData); i++ {
		value, err := t.readTimeSeriesAdjustedItem(csvData[i])
		if err != nil {
			return nil, err
		}

		value.Symbol = p.Symbol
		result = append(result, value)
	}

	return result, nil
}

func (t *timeSeriesClient) createRequestUrl(p timeSeriesParameter) string {
	endpoint := &url.URL{}
	endpoint.Scheme = _Alphavantage_Http_Scheme
	endpoint.Host = _Alphavantage_Host
	endpoint.Path = _Alphavantage_Path
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

func (t *timeSeriesClient) readTimeSeriesItem(s []string) (*TimeSeries, error) {
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

func (t *timeSeriesClient) readTimeSeriesAdjustedItem(s []string) (*TimeSeriesAdjusted, error) {
	const (
		timestamp = iota
		open
		high
		low
		close
		adjusted_close
		volume
		dividend_amount
		split_coefficient
	)

	value := &TimeSeriesAdjusted{}

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

	f, err = parseFloat(s[adjusted_close])
	if err != nil {
		return nil, fmt.Errorf("error parsing adjusted_close %s", s[adjusted_close])
	}
	value.AdjustedClose = f

	f, err = parseFloat(s[dividend_amount])
	if err != nil {
		return nil, fmt.Errorf("error parsing dividend_amount %s", s[dividend_amount])
	}
	value.DividendAmount = f

	// this for TIME_SERIES_DAILY_ADJUSTED API, and it is Premium.
	if len(s) > split_coefficient {
		f, err = parseFloat(s[split_coefficient])
		if err != nil {
			return nil, fmt.Errorf("error parsing split_coefficient %s", s[split_coefficient])
		}
		value.SplitCoefficient = f
	}

	return value, nil
}

func (t *timeSeriesClient) checkTimeSeriesParamter(p timeSeriesParameter) error {
	if len(strings.Trim(p.Symbol, " ")) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}
	_, err := t.checkTimeSeriesFunction(timeSeriesDataFunctionList, p.Function)
	if err != nil {
		_, err := t.checkTimeSeriesFunction(timeSeriesDataaAjustedFunctionList, p.Function)
		if err != nil {
			return err
		}
	}

	if p.Function == _TIME_SERIES_INTRADAY {
		return t.checkTimeSeriesInterval(p.Interval)
	}

	return nil
}

func (t *timeSeriesClient) checkTimeSeriesAdjustedParamter(p timeSeriesParameter) error {
	if len(strings.Trim(p.Symbol, " ")) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}
	_, err := t.checkTimeSeriesFunction(timeSeriesDataaAjustedFunctionList, p.Function)
	if err != nil {
		return err
	}

	if p.Function == _TIME_SERIES_INTRADAY {
		return t.checkTimeSeriesInterval(p.Interval)
	}

	return nil
}

func (t *timeSeriesClient) checkTimeSeriesFunction(functionList []string, function string) (*timeSeriesFunctionType, error) {
	if have := slices.Contains(functionList, function); !have {
		return nil, fmt.Errorf("invalid function name %s", function)
	}

	result := &timeSeriesFunctionType{
		name: function,
	}
	return result, nil
}

func (t *timeSeriesClient) checkTimeSeriesInterval(interval string) error {
	if have := slices.Contains(timeSeriesDataIntervalList, interval); !have {
		return fmt.Errorf("invalid interval name %s", interval)
	}
	return nil
}
