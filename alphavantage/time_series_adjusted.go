package alphavantage

import (
	"fmt"
	"strings"
)

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

// This API returns raw (as-traded) daily open/high/low/close/volume values, adjusted close values, and historical split/dividend events of the global equity specified,
// covering 20+ years of historical data. The OHLCV data is sometimes called "candles" in finance literature.
func (t *TimeSeriesClient) TimeSeriesDailyAdjusted(p TimeSeriesDailyAdjustedParameter) ([]*TimeSeriesAdjusted, error) {
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

// This API returns weekly adjusted time series
// (last trading day of each week, weekly open, weekly high, weekly low, weekly close, weekly adjusted close, weekly volume, weekly dividend)
// of the global equity specified, covering 20+ years of historical data.
func (t *TimeSeriesClient) TimeSeriesWeeklyAdjusted(p TimeSeriesWeeklyAdjustedParameter) ([]*TimeSeriesAdjusted, error) {
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

// This API returns monthly adjusted time series 
// (last trading day of each month, monthly open, monthly high, monthly low, monthly close, monthly adjusted close, monthly volume, monthly dividend)
// of the equity specified, covering 20+ years of historical data.
func (t *TimeSeriesClient) TimeSeriesMonthlyAdjusted(p TimeSeriesMonthlyAdjustedParameter) ([]*TimeSeriesAdjusted, error) {
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

func (t *TimeSeriesClient) readTimeSeriesAdjusted(p timeSeriesParameter) ([]*TimeSeriesAdjusted, error) {
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

func (t *TimeSeriesClient) readTimeSeriesAdjustedItem(s []string) (*TimeSeriesAdjusted, error) {
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
