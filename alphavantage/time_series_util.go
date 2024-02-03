package alphavantage

import (
	"fmt"
	"strings"
	"time"
)

// symbol,open,high,low,price,volume,latestDay,previousClose,change,changePercent
type GlobalQuote struct {
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

// parameter for GLOBAL_QUOTE API
type GlobalQuoteParameter struct {
	// The name of the equity of your choice. For example: symbol=IBM
	Symbol string `json:"symbol"`
}

func (t GlobalQuoteParameter) Validation() error {
	if len(strings.Trim(t.Symbol, " ")) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}

	return nil
}

func (t *TimeSeriesClient) GlobalQuote(p GlobalQuoteParameter) ([]*GlobalQuote, error) {
	err := p.Validation()
	if err != nil {
		return nil, err
	}

	innnerParameter := timeSeriesParameter{
		Function:   _GLOBAL_QUOTE,
		Symbol:     p.Symbol,
	}

	path := t.createRequestUrl(innnerParameter)
	csvData, err := t.httpClient.getCsv(path)
	if err != nil {
		return nil, err
	}

	result := make([]*GlobalQuote, 0)

	for i := 0; i < len(csvData); i++ {
		value, err := t.readGlobalQuoteItem(csvData[i])
		if err != nil {
			return nil, err
		}

		value.Symbol = p.Symbol
		result = append(result, value)
	}

	return result, nil
}

func (t *TimeSeriesClient) readGlobalQuoteItem(s []string) (*GlobalQuote, error) {
	const (
		symbol = iota
		open
		high
		low
		price
		volume
		latestDay
		previousClose
		change
		changePercent
	)

	value := &GlobalQuote{}
	value.Symbol = s[symbol]

	d, err := parseFloat(s[open])
	if err != nil {
		return nil, fmt.Errorf("error parsing open %s", s[open])
	}
	value.Open = d

	f, err := parseFloat(s[high])
	if err != nil {
		return nil, fmt.Errorf("error parsing high %s", s[high])
	}
	value.High = f

	f, err = parseFloat(s[low])
	if err != nil {
		return nil, fmt.Errorf("error parsing low %s", s[low])
	}
	value.Low = f

	f, err = parseFloat(s[price])
	if err != nil {
		return nil, fmt.Errorf("error parsing price %s", s[price])
	}
	value.Price = f

	f, err = parseFloat(s[volume])
	if err != nil {
		return nil, fmt.Errorf("error parsing volume %s", s[volume])
	}
	value.Volume = f

	ti, err := parseTime(s[latestDay])
	if err != nil {
		return nil, fmt.Errorf("error parsing latestDay %s", s[latestDay])
	}
	value.LatestDay = ti

	f, err = parseFloat(s[previousClose])
	if err != nil {
		return nil, fmt.Errorf("error parsing previous_close %s", s[previousClose])
	}
	value.PreviousClose = f

	f, err = parseFloat(s[change])
	if err != nil {
		return nil, fmt.Errorf("error parsing change %s", s[change])
	}
	value.Change = f

	f, err = parseFloat(s[changePercent])
	if err != nil {
		return nil, fmt.Errorf("error parsing change_percent %s", s[changePercent])
	}
	value.ChangePercent = f

	return value, nil
}
