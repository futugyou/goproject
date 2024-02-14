package alphavantage

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/futugyou/alphavantage/enums"
	"github.com/gopherjs/gopherjs/compiler/natives/src/strconv"
)

type AnalyticsClient struct {
	innerClient
}

// https://www.alphavantage.co/documentation/#analytics-fixed-window
func NewAnalyticsClient(apikey string) *AnalyticsClient {
	return &AnalyticsClient{
		innerClient{
			httpClient: newHttpClient(),
			apikey:     apikey,
		},
	}
}

type AnalyticsFixedWindowParameter struct {
	Symbols      string         `json:"symbols"`
	Rannge       string         `json:"rannge"`
	OHLC         enums.OHLC     `json:"ohlc"`
	Interval     enums.Interval `json:"interval"`
	Calculations string         `json:"calculations"`
}

func (p AnalyticsFixedWindowParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	if len(strings.TrimSpace(p.Symbols)) == 0 {
		return nil, fmt.Errorf("symbol not be empty or whitespace")
	}
	dic["SYMBOLS"] = p.Symbols
	if len(strings.TrimSpace(p.Rannge)) == 0 {
		return nil, fmt.Errorf("rannge not be empty or whitespace")
	}
	dic["RANGE"] = p.Rannge
	if len(strings.TrimSpace(p.Calculations)) == 0 {
		return nil, fmt.Errorf("calculations not be empty or whitespace")
	}
	dic["CALCULATIONS"] = p.Calculations
	if p.Interval == nil {
		return nil, fmt.Errorf("interval not be empty or whitespace")
	}
	if p.Interval != nil {
		dic["INTERVAL"] = p.Interval.String()
	}
	if p.OHLC != nil {
		dic["OHLC"] = p.OHLC.String()
	}

	dic["datatype"] = "json"
	return dic, nil
}

type AnalyticsResult struct {
	MetaData MetaData `json:"meta_data"`
	Payload  Payload  `json:"payload"`
}

type MetaData struct {
	Symbols  string `json:"symbols"`
	MinDt    string `json:"min_dt"`
	MaxDt    string `json:"max_dt"`
	Ohlc     string `json:"ohlc"`
	Interval string `json:"interval"`
}

type Payload struct {
	ReturnsCalculations map[string]CalculationItem `json:"RETURNS_CALCULATIONS"`
}

type CalculationItem map[string]float64

func (t *AnalyticsClient) AnalyticsFixedWindow(p AnalyticsFixedWindowParameter) (*AnalyticsResult, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}
	path := t.createUrl(dic, "timeseries/analytics")
	result := &AnalyticsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type AnalyticsSlidingWindowParameter struct {
	Symbols      string         `json:"symbols"`
	Rannge       string         `json:"rannge"`
	OHLC         enums.OHLC     `json:"ohlc"`
	Interval     enums.Interval `json:"interval"`
	Calculations string         `json:"calculations"`
	WindowSize   int            `json:"window_size"`
}

func (p AnalyticsSlidingWindowParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	if len(strings.TrimSpace(p.Symbols)) == 0 {
		return nil, fmt.Errorf("symbol not be empty or whitespace")
	}
	dic["SYMBOLS"] = p.Symbols
	if len(strings.TrimSpace(p.Rannge)) == 0 {
		return nil, fmt.Errorf("rannge not be empty or whitespace")
	}
	dic["RANGE"] = p.Rannge
	if len(strings.TrimSpace(p.Calculations)) == 0 {
		return nil, fmt.Errorf("calculations not be empty or whitespace")
	}
	dic["CALCULATIONS"] = p.Calculations
	if p.Interval == nil {
		return nil, fmt.Errorf("interval not be empty or whitespace")
	}
	if p.Interval != nil {
		dic["INTERVAL"] = p.Interval.String()
	}
	if p.OHLC != nil {
		dic["OHLC"] = p.OHLC.String()
	}
	if p.WindowSize < 10 {
		p.WindowSize = 10
	}
	dic["WINDOW_SIZE"] = strconv.Itoa(p.WindowSize)
	dic["datatype"] = "json"
	return dic, nil
}

type AnalyticsSlidingResult struct {
	MetaData SlidingMetaData `json:"meta_data"`
	Payload  SlidingPayload  `json:"payload"`
}

type SlidingMetaData struct {
	Symbols    string `json:"symbols"`
	WindowSize int64  `json:"window_size"`
	MinDt      string `json:"min_dt"`
	MaxDt      string `json:"max_dt"`
	Ohlc       string `json:"ohlc"`
	Interval   string `json:"interval"`
}

type SlidingPayload struct {
	ReturnsCalculations map[string]map[string]map[string]map[string]float64 `json:"RETURNS_CALCULATIONS"`
}

func (t *AnalyticsClient) AnalyticsSlidingWindow(p AnalyticsSlidingWindowParameter) (*AnalyticsResult, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}
	path := t.createUrl(dic, "timeseries/running_analytics")
	result := &AnalyticsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
