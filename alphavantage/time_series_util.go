package alphavantage

import (
	"fmt"
	"strings"
	"time"
)

// symbol,open,high,low,price,volume,latestDay,previousClose,change,changePercent
type GlobalQuote struct {
	Symbol        string    `json:"symbol" csv:"symbol"`
	Open          float64   `json:"open" csv:"open"`
	High          float64   `json:"high" csv:"high"`
	Low           float64   `json:"low" csv:"low"`
	Price         float64   `json:"price" csv:"price"`
	Volume        float64   `json:"volume" csv:"volume"`
	LatestDay     time.Time `json:"latestDay" csv:"latestDay"`
	PreviousClose float64   `json:"previous_close" csv:"previousClose"`
	Change        float64   `json:"change" csv:"change"`
	ChangePercent float64   `json:"change_percent" csv:"changePercent"`
}

// symbol,name,type,region,marketOpen,marketClose,timezone,currency,matchScore
type SymbolSearch struct {
	Symbol      string  `json:"symbol" csv:"symbol"`
	Name        string  `json:"name" csv:"name"`
	Type        string  `json:"type" csv:"type"`
	Region      string  `json:"region" csv:"region"`
	MarketOpen  string  `json:"marketOpen" csv:"marketOpen"`
	MarketClose string  `json:"marketClose" csv:"marketClose"`
	Timezone    string  `json:"timezone" csv:"timezone"`
	Currency    string  `json:"currency" csv:"currency"`
	MatchScore  float64 `json:"matchScore" csv:"matchScore"`
}

type MarketStatus struct {
	Endpoint    string   `json:"endpoint"`
	Markets     []Market `json:"markets"`
	Information string   `json:"Information"`
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
	if len(strings.TrimSpace(t.Symbol)) == 0 {
		return fmt.Errorf("symbol can not be empty or whitespace")
	}

	return nil
}

// parameter for SYMBOL_SEARCH API
type SymbolSearchParameter struct {
	// A text string of your choice. For example: keywords=microsoft.
	Keywords string `json:"keywords"`
}

func (t SymbolSearchParameter) Validation() error {
	if len(strings.TrimSpace(t.Keywords)) == 0 {
		return fmt.Errorf("keywords can not be empty or whitespace")
	}

	return nil
}

// A lightweight alternative to the time series APIs, this service returns the latest price and volume information for a ticker of your choice.
func (t *TimeSeriesClient) GlobalQuote(p GlobalQuoteParameter) ([]*GlobalQuote, error) {
	err := p.Validation()
	if err != nil {
		return nil, err
	}

	innnerParameter := timeSeriesParameter{
		Function: _GLOBAL_QUOTE,
		Symbol:   p.Symbol,
	}

	path := t.createRequestUrl(innnerParameter)
	result := make([]*GlobalQuote, 0)

	err = t.httpClient.getCsvByUtil(path, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// A lightweight alternative to the time series APIs, this service returns the latest price and volume information for a ticker of your choice.
func (t *TimeSeriesClient) SymbolSearch(p SymbolSearchParameter) ([]*SymbolSearch, error) {
	err := p.Validation()
	if err != nil {
		return nil, err
	}

	innnerParameter := timeSeriesParameter{
		Function:   _SYMBOL_SEARCH,
		Dictionary: map[string]string{"keywords": p.Keywords},
	}

	path := t.createRequestUrl(innnerParameter)
	result := make([]*SymbolSearch, 0)

	err = t.httpClient.getCsvByUtil(path, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// This endpoint returns the current market status (open vs. closed) of major trading venues for equities, forex, and cryptocurrencies around the world.
func (t *TimeSeriesClient) MarketStatus() (*MarketStatus, error) {
	innnerParameter := timeSeriesParameter{
		Function: _MARKET_STATUS,
	}

	path := t.createRequestUrl(innnerParameter)
	result := &MarketStatus{}
	err := t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
