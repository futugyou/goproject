package alphavantage

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type DigitalCurrencyClient struct {
	innerClient
}

// APIs under this section provide a wide range of data feed for digital and crypto currencies such as Bitcoin.
func NewDigitalCurrencyClient(apikey string) *DigitalCurrencyClient {
	return &DigitalCurrencyClient{
		innerClient{
			httpClient: newHttpClient(),
			apikey:     apikey,
		},
	}
}

// parameter for CURRENCY_EXCHANGE_RATE API
type CryptoExchangeParameter struct {
	// The currency you would like to get the exchange rate for. It can either be a physical currency or digital/crypto currency.
	// For example: from_currency=USD or from_currency=BTC.
	FromCurrency string `json:"from_currency"`
	// The destination currency for the exchange rate. It can either be a physical currency or digital/crypto currency.
	// For example: to_currency=USD or to_currency=BTC.
	ToCurrency string `json:"to_currency"`
}

func (p CryptoExchangeParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = "CURRENCY_EXCHANGE_RATE"
	if len(strings.TrimSpace(p.FromCurrency)) == 0 {
		return nil, fmt.Errorf("from_symbol not be empty or whitespace")
	}
	dic["from_currency"] = strings.TrimSpace(p.FromCurrency)

	if len(strings.TrimSpace(p.ToCurrency)) == 0 {
		return nil, fmt.Errorf("to_symbol not be empty or whitespace")
	}
	dic["to_currency"] = strings.TrimSpace(p.ToCurrency)

	dic["datatype"] = "csv"
	return dic, nil
}

type CryptoExchang struct {
	DigitalExchangeRate DigitalExchangeRate `json:"Realtime Currency Exchange Rate"`
}

type DigitalExchangeRate struct {
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

func (t *DigitalCurrencyClient) CryptoExchange(p CryptoExchangeParameter) (*CryptoExchang, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &CryptoExchang{}

	err = t.httpClient.getJson(path, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// parameter for CRYPTO_INTRADAY API
type CryptoIntradayParameter struct {
	// The digital/crypto currency of your choice. It can be any of the currencies in the digital currency list. For example: symbol=ETH.
	Symbol string `json:"symbol"`
	// The exchange market of your choice. It can be any of the market in the market list. For example: market=USD.
	Market string `json:"market"`
	// Time interval between two consecutive data points in the time series. The following values are supported: 1min, 5min, 15min, 30min, 60min
	Interval string `json:"interval"`
}

func (p CryptoIntradayParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = "CRYPTO_INTRADAY"
	if len(strings.TrimSpace(p.Symbol)) == 0 {
		return nil, fmt.Errorf("symbol not be empty or whitespace")
	}
	dic["symbol"] = strings.TrimSpace(p.Symbol)

	if len(strings.TrimSpace(p.Market)) == 0 {
		return nil, fmt.Errorf("market not be empty or whitespace")
	}
	dic["market"] = strings.TrimSpace(p.Market)

	if slices.Contains(timeSeriesDataIntervalList, strings.TrimSpace(p.Interval)) {
		dic["interval"] = strings.TrimSpace(p.Interval)
	} else {
		return nil, fmt.Errorf("interval only can be %s", strings.Join(timeSeriesDataIntervalList, ","))
	}

	dic["datatype"] = "csv"
	return dic, nil
}

// timestamp,open,high,low,close
type CryptoIntraday struct {
	Symbol    string    `json:"symbol" csv:"-"`
	Market    string    `json:"market" csv:"-"`
	Timestamp time.Time `json:"timestamp" csv:"timestamp"`
	Open      float64   `json:"open" csv:"open"`
	High      float64   `json:"high" csv:"high"`
	Low       float64   `json:"low" csv:"low"`
	Close     float64   `json:"close" csv:"close"`
	Volume    float64   `json:"volume" csv:"volume"`
}

func (t *DigitalCurrencyClient) CryptoIntraday(p CryptoIntradayParameter) ([]CryptoIntraday, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := make([]CryptoIntraday, 0)

	err = t.httpClient.getCsvByUtil(path, &result)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result); i++ {
		result[i].Symbol = p.Symbol
		result[i].Market = p.Market
	}

	return result, nil
}
