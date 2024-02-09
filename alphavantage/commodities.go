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

type innerCommoditiesParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.LongInterval `json:"interval"`
}

func (p innerCommoditiesParameter) Validation(function string) (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = function
	if p.Interval != nil {
		dic["interval"] = p.Interval.String()
	}

	dic["datatype"] = "json"
	return dic, nil
}

type innerCommoditiesResult struct {
	Name     string  `json:"name"`
	Interval string  `json:"interval"`
	Unit     string  `json:"unit"`
	Data     []Datum `json:"data"`
}

type Datum struct {
	Date  string `json:"date"`
	Value string `json:"value"`
}

func (t *CommoditiesClient) innerCommoditiesRequest(p innerCommoditiesParameter, function string) (*innerCommoditiesResult, error) {
	dic, err := p.Validation(function)
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &innerCommoditiesResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// parameter for WTI API
type CrudeOilWtiParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.LongInterval `json:"interval"`
}

type CrudeOilWti struct {
	innerCommoditiesResult
}

// This API returns the West Texas Intermediate (WTI) crude oil prices in daily, weekly, and monthly horizons.
// Source: U.S. Energy Information Administration, Crude Oil Prices: West Texas Intermediate (WTI) -
// Cushing, Oklahoma, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) CrudeOilWti(p CrudeOilWtiParameter) (*CrudeOilWti, error) {
	pp := innerCommoditiesParameter(p)
	inner, err := t.innerCommoditiesRequest(pp, "WTI")

	if err != nil {
		return nil, err
	}

	result := &CrudeOilWti{*inner}
	return result, nil
}

// parameter for BRENT API
type CrudeOilBrentParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.LongInterval `json:"interval"`
}

type CrudeOilBrent struct {
	innerCommoditiesResult
}

// This API returns the Brent (Europe) crude oil prices in daily, weekly, and monthly horizons.
// Source: U.S. Energy Information Administration, Crude Oil Prices: Brent - Europe, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) CrudeOilBrent(p CrudeOilBrentParameter) (*CrudeOilBrent, error) {
	pp := innerCommoditiesParameter(p)
	inner, err := t.innerCommoditiesRequest(pp, "BRENT")

	if err != nil {
		return nil, err
	}

	result := &CrudeOilBrent{*inner}
	return result, nil
}
