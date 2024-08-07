package alphavantage

import (
	"fmt"
	"strings"
	"time"

	"github.com/futugyou/alphavantage/enums"
)

type ForeignExchangeRatesClient struct {
	innerClient
}

// APIs under this section provide a wide range of data feed for realtime and historical forex (FX) rates.
func NewForeignExchangeRatesClient(apikey string) *ForeignExchangeRatesClient {
	return &ForeignExchangeRatesClient{
		innerClient{
			httpClient: newHttpClient(),
			apikey:     apikey,
		},
	}
}

// parameter for CURRENCY_EXCHANGE_RATE API
type CurrencyExchangeParameter struct {
	// The currency you would like to get the exchange rate for. It can either be a physical currency or digital/crypto currency.
	// For example: from_currency=USD or from_currency=BTC.
	FromCurrency string `json:"from_currency"`
	// The destination currency for the exchange rate. It can either be a physical currency or digital/crypto currency.
	// For example: to_currency=USD or to_currency=BTC.
	ToCurrency string `json:"to_currency"`
}

type CurrencyExchange struct {
	ForeignExchangeRate ForeignExchangeRate `json:"Realtime Currency Exchange Rate"`
	Information         string              `json:"Information"`
	ErrorMessage        string              `json:"Error Message"`
}

type ForeignExchangeRate struct {
	FromCurrencyCode string `json:"1. From_Currency Code"`
	FromCurrencyName string `json:"2. From_Currency Name"`
	ToCurrencyCode   string `json:"3. To_Currency Code"`
	ToCurrencyName   string `json:"4. To_Currency Name"`
	ExchangeRate     string `json:"5. Exchange Rate"`
	LastRefreshed    string `json:"6. Last Refreshed"`
	TimeZone         string `json:"7. Time Zone"`
	BidPrice         string `json:"8. Bid Price"`
	AskPrice         string `json:"9. Ask Price"`
}

// This API returns the realtime exchange rate for a pair of digital currency (e.g., Bitcoin) and physical currency (e.g., USD).
func (t *ForeignExchangeRatesClient) CurrencyExchange(p CurrencyExchangeParameter) (*CurrencyExchange, error) {
	dic := make(map[string]string)
	dic["function"] = "CURRENCY_EXCHANGE_RATE"
	dic["from_currency"] = p.FromCurrency
	dic["to_currency"] = p.ToCurrency

	path := t.createQuerytUrl(dic)
	result := &CurrencyExchange{}

	err := t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// parameter for FX_INTRADAY API
type FxIntradayParameter struct {
	// A three-letter symbol from the forex currency list. For example: from_symbol=EUR
	FromSymbol string `json:"from_symbol"`
	// A three-letter symbol from the forex currency list. For example: from_symbol=EUR
	ToSymbol string `json:"to_symbol"`
	// Time interval between two consecutive data points in the time series. The following values are supported: 1min, 5min, 15min, 30min, 60min
	Interval enums.TimeInterval `json:"interval"`
}

func (p FxIntradayParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = "FX_INTRADAY"
	if len(strings.TrimSpace(p.FromSymbol)) == 0 {
		return nil, fmt.Errorf("from_symbol not be empty or whitespace")
	}
	dic["from_symbol"] = strings.TrimSpace(p.FromSymbol)

	if len(strings.TrimSpace(p.ToSymbol)) == 0 {
		return nil, fmt.Errorf("to_symbol not be empty or whitespace")
	}
	dic["to_symbol"] = strings.TrimSpace(p.ToSymbol)

	dic["interval"] = p.Interval.String()

	dic["datatype"] = "csv"
	return dic, nil
}

// timestamp,open,high,low,close
type FxIntraday struct {
	FromSymbol string    `json:"from_symbol" csv:"-"`
	ToSymbol   string    `json:"to_symbol" csv:"-"`
	Timestamp  time.Time `json:"timestamp" csv:"timestamp"`
	Open       float64   `json:"open" csv:"open"`
	High       float64   `json:"high" csv:"high"`
	Low        float64   `json:"low" csv:"low"`
	Close      float64   `json:"close" csv:"close"`
}

// This API returns intraday time series (timestamp, open, high, low, close) of the FX currency pair specified, updated realtime.
func (t *ForeignExchangeRatesClient) FxIntraday(p FxIntradayParameter) ([]FxIntraday, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := make([]FxIntraday, 0)

	err = t.httpClient.getCsvByUtil(path, &result)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result); i++ {
		result[i].FromSymbol = p.FromSymbol
		result[i].ToSymbol = p.ToSymbol
	}

	return result, nil
}

// parameter for FX_DAILY API
type FxDailyParameter struct {
	// A three-letter symbol from the forex currency list. For example: from_symbol=EUR
	FromSymbol string `json:"from_symbol"`
	// A three-letter symbol from the forex currency list. For example: from_symbol=EUR
	ToSymbol string `json:"to_symbol"`
}

func (p FxDailyParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = "FX_DAILY"
	if len(strings.TrimSpace(p.FromSymbol)) == 0 {
		return nil, fmt.Errorf("from_symbol not be empty or whitespace")
	}
	dic["from_symbol"] = strings.TrimSpace(p.FromSymbol)

	if len(strings.TrimSpace(p.ToSymbol)) == 0 {
		return nil, fmt.Errorf("to_symbol not be empty or whitespace")
	}
	dic["to_symbol"] = strings.TrimSpace(p.ToSymbol)

	dic["datatype"] = "csv"
	return dic, nil
}

// timestamp,open,high,low,close
type FxDaily struct {
	FromSymbol string    `json:"from_symbol" csv:"-"`
	ToSymbol   string    `json:"to_symbol" csv:"-"`
	Timestamp  time.Time `json:"timestamp" csv:"timestamp"`
	Open       float64   `json:"open" csv:"open"`
	High       float64   `json:"high" csv:"high"`
	Low        float64   `json:"low" csv:"low"`
	Close      float64   `json:"close" csv:"close"`
}

// This API returns the daily time series (timestamp, open, high, low, close) of the FX currency pair specified, updated realtime.
func (t *ForeignExchangeRatesClient) FxDaily(p FxDailyParameter) ([]FxDaily, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := make([]FxDaily, 0)

	err = t.httpClient.getCsvByUtil(path, &result)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result); i++ {
		result[i].FromSymbol = p.FromSymbol
		result[i].ToSymbol = p.ToSymbol
	}

	return result, nil
}

// parameter for FX_WEEKLY API
type FxWeeklyParameter struct {
	// A three-letter symbol from the forex currency list. For example: from_symbol=EUR
	FromSymbol string `json:"from_symbol"`
	// A three-letter symbol from the forex currency list. For example: from_symbol=EUR
	ToSymbol string `json:"to_symbol"`
}

func (p FxWeeklyParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = "FX_WEEKLY"
	if len(strings.TrimSpace(p.FromSymbol)) == 0 {
		return nil, fmt.Errorf("from_symbol not be empty or whitespace")
	}
	dic["from_symbol"] = strings.TrimSpace(p.FromSymbol)

	if len(strings.TrimSpace(p.ToSymbol)) == 0 {
		return nil, fmt.Errorf("to_symbol not be empty or whitespace")
	}
	dic["to_symbol"] = strings.TrimSpace(p.ToSymbol)

	dic["datatype"] = "csv"
	return dic, nil
}

// timestamp,open,high,low,close
type FxWeekly struct {
	FromSymbol string    `json:"from_symbol" csv:"-"`
	ToSymbol   string    `json:"to_symbol" csv:"-"`
	Timestamp  time.Time `json:"timestamp" csv:"timestamp"`
	Open       float64   `json:"open" csv:"open"`
	High       float64   `json:"high" csv:"high"`
	Low        float64   `json:"low" csv:"low"`
	Close      float64   `json:"close" csv:"close"`
}

// This API returns the weekly time series (timestamp, open, high, low, close) of the FX currency pair specified, updated realtime.
// The latest data point is the price information for the week (or partial week) containing the current trading day, updated realtime.
func (t *ForeignExchangeRatesClient) FxWeekly(p FxWeeklyParameter) ([]FxWeekly, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := make([]FxWeekly, 0)

	err = t.httpClient.getCsvByUtil(path, &result)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result); i++ {
		result[i].FromSymbol = p.FromSymbol
		result[i].ToSymbol = p.ToSymbol
	}

	return result, nil
}

// parameter for FX_MONTHLY API
type FxMonthlyParameter struct {
	// A three-letter symbol from the forex currency list. For example: from_symbol=EUR
	FromSymbol string `json:"from_symbol"`
	// A three-letter symbol from the forex currency list. For example: from_symbol=EUR
	ToSymbol string `json:"to_symbol"`
}

func (p FxMonthlyParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = "FX_MONTHLY"
	if len(strings.TrimSpace(p.FromSymbol)) == 0 {
		return nil, fmt.Errorf("from_symbol not be empty or whitespace")
	}
	dic["from_symbol"] = strings.TrimSpace(p.FromSymbol)

	if len(strings.TrimSpace(p.ToSymbol)) == 0 {
		return nil, fmt.Errorf("to_symbol not be empty or whitespace")
	}
	dic["to_symbol"] = strings.TrimSpace(p.ToSymbol)

	dic["datatype"] = "csv"
	return dic, nil
}

// timestamp,open,high,low,close
type FxMonthly struct {
	FromSymbol string    `json:"from_symbol" csv:"-"`
	ToSymbol   string    `json:"to_symbol" csv:"-"`
	Timestamp  time.Time `json:"timestamp" csv:"timestamp"`
	Open       float64   `json:"open" csv:"open"`
	High       float64   `json:"high" csv:"high"`
	Low        float64   `json:"low" csv:"low"`
	Close      float64   `json:"close" csv:"close"`
}

// This API returns the monthly time series (timestamp, open, high, low, close) of the FX currency pair specified, updated realtime.
// The latest data point is the prices information for the month (or partial month) containing the current trading day, updated realtime.
func (t *ForeignExchangeRatesClient) FxMonthly(p FxMonthlyParameter) ([]FxMonthly, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := make([]FxMonthly, 0)

	err = t.httpClient.getCsvByUtil(path, &result)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result); i++ {
		result[i].FromSymbol = p.FromSymbol
		result[i].ToSymbol = p.ToSymbol
	}

	return result, nil
}
