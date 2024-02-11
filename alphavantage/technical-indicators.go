package alphavantage

import (
	"fmt"
	"strings"
	"time"

	"github.com/futugyou/alphavantage/enums"
	"github.com/futugyou/alphavantage/functions"
)

type TechnicalIndicatorsClient struct {
	innerClient
}

// Technical indicator APIs for a given equity or currency exchange pair, derived from the underlying time series based stock API and forex data.
// All indicators are calculated from adjusted time series data to eliminate artificial price/volume perturbations from historical split and dividend events.
// IMPORTANT!!! - This category has 53 APIs, so the request parameters will not be checked at compile time.
// Please refer to the documentation for details.
// https://www.alphavantage.co/documentation/#technical-indicators
func NewTechnicalIndicatorsClient(apikey string) *TechnicalIndicatorsClient {
	return &TechnicalIndicatorsClient{
		innerClient{
			httpClient: newHttpClient(),
			apikey:     apikey,
		},
	}
}

type TechnicalIndicatorsParameter struct {
	Function functions.TechnicalIndicatorsType `json:"function"`
	// The name of the ticker of your choice. For example: symbol=IBM
	Symbol string `json:"symbol"`
	// Interval   enums.LongInterval `json:"interval"`
	Month      string            `json:"month"`
	TimePeriod float64           `json:"time_period"`
	SeriesType enums.SeriesType  `json:"series_type"`
	Interval   enums.Interval    `json:"interval"`
	Dictionary map[string]string `json:"dictionary"`
}

func (p TechnicalIndicatorsParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = p.Function.String()
	dic["datatype"] = "csv"
	if len(strings.TrimSpace(p.Symbol)) == 0 {
		return nil, fmt.Errorf("symbol can not be null or whtiespace")
	}
	dic["symbol"] = strings.TrimSpace(p.Symbol)
	if len(strings.TrimSpace(p.Month)) != 0 {
		dic["month"] = strings.TrimSpace(p.Month)
	}
	if p.SeriesType != nil {
		dic["series_type"] = p.SeriesType.String()
	}
	if p.Interval != nil {
		dic["interval"] = p.Interval.String()
	}

	for k, v := range p.Dictionary {
		dic[k] = v
	}

	return dic, nil
}

type TechnicalIndicatorsResult struct {
	Symbol   string    `json:"symbol"`
	Function string    `json:"function"`
	Time     time.Time `json:"time"`
	Item1    float64   `json:"item1"`
	Item2    float64   `json:"item2"`
	Item3    float64   `json:"item3"`
}

func (t *TechnicalIndicatorsClient) GetTechnicalIndicatorsData(p TechnicalIndicatorsParameter) ([]TechnicalIndicatorsResult, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	csvData, err := t.httpClient.getCsv(path)
	if err != nil {
		return nil, err
	}

	result := make([]TechnicalIndicatorsResult, 0)

	for i := 0; i < len(csvData); i++ {
		value, err := t.readTechnicalIndicatorsItem(csvData[i])
		if err != nil {
			return nil, err
		}

		value.Symbol = p.Symbol
		value.Function = p.Function.String()
		result = append(result, *value)
	}

	return result, nil
}

func (t *TechnicalIndicatorsClient) readTechnicalIndicatorsItem(s []string) (*TechnicalIndicatorsResult, error) {
	const (
		timestamp = iota
		item1
		item2
		item3
	)

	value := &TechnicalIndicatorsResult{}

	d, err := parseTime(s[timestamp])
	if err != nil {
		return nil, fmt.Errorf("error parsing timestamp %s", s[timestamp])
	}
	value.Time = d

	f, err := parseFloat(s[item1])
	if err != nil {
		return nil, fmt.Errorf("error parsing item1 %s", s[item1])
	}
	value.Item1 = f

	if item2 < len(s) {
		f, err := parseFloat(s[item2])
		if err != nil {
			return nil, fmt.Errorf("error parsing item2 %s", s[item2])
		}
		value.Item2 = f
	}

	if item3 < len(s) {
		f, err := parseFloat(s[item3])
		if err != nil {
			return nil, fmt.Errorf("error parsing item3 %s", s[item3])
		}
		value.Item3 = f
	}

	return value, nil
}
