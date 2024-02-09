package alphavantage

import (
	"github.com/futugyou/alphavantage/enums"
)

type CommoditiesClient struct {
	innerClient
}

// APIs under this section provide historical data for major commodities such as crude oil, natural gas, copper, wheat, etc.,
// spanning across various temporal horizons (daily, weekly, monthly, quarterly, etc.)
func NewCommoditiesClient(apikey string) *CommoditiesClient {
	return &CommoditiesClient{
		innerClient{
			httpClient: newHttpClient(),
			apikey:     apikey,
		},
	}
}

// parameter for WTI API
type CrudeOilWtiParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.LongInterval `json:"interval"`
}

func (p CrudeOilWtiParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = "WTI"
	if p.Interval != nil {
		dic["interval"] = p.Interval.String()
	}

	dic["datatype"] = "json"
	return dic, nil
}

type CrudeOilWti struct {
	Name     string  `json:"name"`
	Interval string  `json:"interval"`
	Unit     string  `json:"unit"`
	Data     []Datum `json:"data"`
}

type Datum struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

// This API returns the West Texas Intermediate (WTI) crude oil prices in daily, weekly, and monthly horizons.
// Source: U.S. Energy Information Administration, Crude Oil Prices: West Texas Intermediate (WTI) -
// Cushing, Oklahoma, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) CrudeOilWti(p CrudeOilWtiParameter) (*CrudeOilWti, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &CrudeOilWti{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
