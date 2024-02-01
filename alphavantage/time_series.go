package alphavantage

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"slices"
	"time"
)

// TODO: update model
// timestamp,open,high,low,close,volume
type TimeSeries struct {
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

type TimeSeriesParameter struct {
	Function   string            `json:"function"`
	Symbol     string            `json:"symbol"`
	Interval   string            `json:"interval"`
	Dictionary map[string]string `json:"dictionary"`
}

type timeSeriesClient struct {
	httpClient *httpClient
	apikey     string
	datatype   string
}

func NewTimeSeriesClient(apikey string) *timeSeriesClient {
	return &timeSeriesClient{
		httpClient: newHttpClient(),
		apikey:     apikey,
		datatype:   _Alphavantage_Datatype,
	}
}

func (t *timeSeriesClient) createRequestUrl(p TimeSeriesParameter) string {
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

func (t *timeSeriesClient) ReadTimeSeries(p TimeSeriesParameter) ([]*TimeSeries, error) {
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

		var value *TimeSeries
		switch p.Function {
		case _TIME_SERIES_INTRADAY, _TIME_SERIES_DAILY, _TIME_SERIES_WEEKLY, _TIME_SERIES_MONTHLY:
			value, err = t.readTimeSeriesItem(record)
		case _TIME_SERIES_DAILY_ADJUSTED, _TIME_SERIES_WEEKLY_ADJUSTED, _TIME_SERIES_MONTHLY_ADJUSTED:
			value, err = t.readTimeSeriesAdjustedItem(record)
		}
		if err != nil {
			return nil, err
		}

		value.Symbol = p.Symbol
		result = append(result, value)
	}

	return result, nil

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

func (t *timeSeriesClient) readTimeSeriesAdjustedItem(s []string) (*TimeSeries, error) {
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

func (t *timeSeriesClient) checkTimeSeriesParameter(function string) (*timeSeriesFunctionType, error) {
	if have := slices.Contains(functionList, function); !have {
		return nil, fmt.Errorf("invalid function name %s", function)
	}

	result := &timeSeriesFunctionType{
		name: function,
	}
	return result, nil
}

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

var functionList = []string{_TIME_SERIES_INTRADAY, _TIME_SERIES_DAILY, _TIME_SERIES_DAILY_ADJUSTED,
	_TIME_SERIES_WEEKLY, _TIME_SERIES_WEEKLY_ADJUSTED, _TIME_SERIES_MONTHLY, _TIME_SERIES_MONTHLY_ADJUSTED,
	_GLOBAL_QUOTE, _SYMBOL_SEARCH, _MARKET_STATUS,
}
