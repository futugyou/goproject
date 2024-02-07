package alphavantage

import (
	"fmt"
	"slices"
	"strings"
	"time"
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
	RealtimeCurrencyExchangeRate RealtimeCurrencyExchangeRate `json:"Realtime Currency Exchange Rate"`
}

type RealtimeCurrencyExchangeRate struct {
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
	Interval string `json:"interval"`
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

	if slices.Contains(timeSeriesDataIntervalList, strings.TrimSpace(p.Interval)) {
		dic["interval"] = strings.TrimSpace(p.Interval)
	} else {
		return nil, fmt.Errorf("interval only can be %s", strings.Join(timeSeriesDataIntervalList, ","))
	}

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
