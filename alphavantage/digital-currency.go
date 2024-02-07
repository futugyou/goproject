package alphavantage

import (
	"fmt"
	"strings"
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
