package alphavantage

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
