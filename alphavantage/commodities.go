package alphavantage

import (
	"fmt"

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
	Interval enums.CommoditiesInterval
	// By default, interval=monthly. Strings monthly, quarterly, and annual are accepted.
	Interval2 enums.CommoditiesInterval2
}

func (p innerCommoditiesParameter) Validation(function string) (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = function
	if p.Interval != nil {
		dic["interval"] = p.Interval.String()
	}
	if p.Interval2 != nil {
		dic["interval"] = p.Interval2.String()
	}

	dic["datatype"] = "json"
	return dic, nil
}

type innerCommoditiesResult struct {
	Name        string  `json:"name"`
	Interval    string  `json:"interval"`
	Unit        string  `json:"unit"`
	Data        []Datum `json:"data"`
	Information string  `json:"Information"`
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

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

// parameter for WTI API
type CrudeOilWtiParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.CommoditiesInterval `json:"interval"`
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
	pp := innerCommoditiesParameter{Interval: p.Interval}
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
	Interval enums.CommoditiesInterval `json:"interval"`
}

type CrudeOilBrent struct {
	innerCommoditiesResult
}

// This API returns the Brent (Europe) crude oil prices in daily, weekly, and monthly horizons.
// Source: U.S. Energy Information Administration, Crude Oil Prices: Brent - Europe, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) CrudeOilBrent(p CrudeOilBrentParameter) (*CrudeOilBrent, error) {
	pp := innerCommoditiesParameter{Interval: p.Interval}
	inner, err := t.innerCommoditiesRequest(pp, "BRENT")

	if err != nil {
		return nil, err
	}

	result := &CrudeOilBrent{*inner}
	return result, nil
}

// parameter for NATURAL_GAS API
type NaturalGasParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.CommoditiesInterval `json:"interval"`
}

type NaturalGas struct {
	innerCommoditiesResult
}

// This API returns the Henry Hub natural gas spot prices in daily, weekly, and monthly horizons.
// Source: U.S. Energy Information Administration, Henry Hub Natural Gas Spot Price, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) NaturalGas(p NaturalGasParameter) (*NaturalGas, error) {
	pp := innerCommoditiesParameter{Interval: p.Interval}
	inner, err := t.innerCommoditiesRequest(pp, "NATURAL_GAS")

	if err != nil {
		return nil, err
	}

	result := &NaturalGas{*inner}
	return result, nil
}

// parameter for COPPER API
type CopperParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.CommoditiesInterval2 `json:"interval"`
}

type Copper struct {
	innerCommoditiesResult
}

// This API returns the global price of copper in monthly, quarterly, and annual horizons.
// Source: International Monetary Fund (IMF Terms of Use), Global price of Copper, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) Copper(p CopperParameter) (*Copper, error) {
	pp := innerCommoditiesParameter{Interval2: p.Interval}
	inner, err := t.innerCommoditiesRequest(pp, "COPPER")

	if err != nil {
		return nil, err
	}

	result := &Copper{*inner}
	return result, nil
}

// parameter for ALUMINUM API
type AluminumParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.CommoditiesInterval2 `json:"interval"`
}

type Aluminum struct {
	innerCommoditiesResult
}

// This API returns the global price of aluminum in monthly, quarterly, and annual horizons.
// Source: International Monetary Fund (IMF Terms of Use), Global price of Aluminum, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) Aluminum(p AluminumParameter) (*Aluminum, error) {
	pp := innerCommoditiesParameter{Interval2: p.Interval}
	inner, err := t.innerCommoditiesRequest(pp, "ALUMINUM")

	if err != nil {
		return nil, err
	}

	result := &Aluminum{*inner}
	return result, nil
}

// parameter for WHEAT API
type WheatParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.CommoditiesInterval2 `json:"interval"`
}

type Wheat struct {
	innerCommoditiesResult
}

// This API returns the global price of wheat in monthly, quarterly, and annual horizons.
// Source: International Monetary Fund (IMF Terms of Use), Global price of Wheat, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) Wheat(p WheatParameter) (*Wheat, error) {
	pp := innerCommoditiesParameter{Interval2: p.Interval}
	inner, err := t.innerCommoditiesRequest(pp, "WHEAT")

	if err != nil {
		return nil, err
	}

	result := &Wheat{*inner}
	return result, nil
}

// parameter for CORN API
type CornParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.CommoditiesInterval2 `json:"interval"`
}

type Corn struct {
	innerCommoditiesResult
}

// This API returns the global price of corn in monthly, quarterly, and annual horizons.
// Source: International Monetary Fund (IMF Terms of Use), Global price of Corn, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) Corn(p CornParameter) (*Corn, error) {
	pp := innerCommoditiesParameter{Interval2: p.Interval}
	inner, err := t.innerCommoditiesRequest(pp, "CORN")

	if err != nil {
		return nil, err
	}

	result := &Corn{*inner}
	return result, nil
}

// parameter for COTTON API
type CottonParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.CommoditiesInterval2 `json:"interval"`
}

type Cotton struct {
	innerCommoditiesResult
}

// This API returns the global price of cotton in monthly, quarterly, and annual horizons.
// Source: International Monetary Fund (IMF Terms of Use), Global price of Cotton, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) Cotton(p CottonParameter) (*Cotton, error) {
	pp := innerCommoditiesParameter{Interval2: p.Interval}
	inner, err := t.innerCommoditiesRequest(pp, "COTTON")

	if err != nil {
		return nil, err
	}

	result := &Cotton{*inner}
	return result, nil
}

// parameter for SUGAR API
type SugarParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.CommoditiesInterval2 `json:"interval"`
}

type Sugar struct {
	innerCommoditiesResult
}

// This API returns the global price of sugar in monthly, quarterly, and annual horizons.
// Source: International Monetary Fund (IMF Terms of Use), Global price of Sugar, No. 11, World, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) Sugar(p SugarParameter) (*Sugar, error) {
	pp := innerCommoditiesParameter{Interval2: p.Interval}
	inner, err := t.innerCommoditiesRequest(pp, "SUGAR")

	if err != nil {
		return nil, err
	}

	result := &Sugar{*inner}
	return result, nil
}

// parameter for COFFEE API
type CoffeeParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.CommoditiesInterval2 `json:"interval"`
}

type Coffee struct {
	innerCommoditiesResult
}

// This API returns the global price of coffee in monthly, quarterly, and annual horizons.
// Source: International Monetary Fund (IMF Terms of Use), Global price of Coffee, Other Mild Arabica, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) Coffee(p CoffeeParameter) (*Coffee, error) {
	pp := innerCommoditiesParameter{Interval2: p.Interval}
	inner, err := t.innerCommoditiesRequest(pp, "COFFEE")

	if err != nil {
		return nil, err
	}

	result := &Coffee{*inner}
	return result, nil
}

// parameter for ALL_COMMODITIES API
type AllCommoditiesParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.CommoditiesInterval2 `json:"interval"`
}

type AllCommodities struct {
	innerCommoditiesResult
}

// This API returns the global price index of all commodities in monthly, quarterly, and annual temporal dimensions.
// Source: International Monetary Fund (IMF Terms of Use), Global Price Index of All Commodities, retrieved from FRED, Federal Reserve Bank of St. Louis.
// This data feed uses the FRED® API but is not endorsed or certified by the Federal Reserve Bank of St. Louis.
// By using this data feed, you agree to be bound by the FRED® API Terms of Use.
func (t *CommoditiesClient) AllCommodities(p AllCommoditiesParameter) (*AllCommodities, error) {
	pp := innerCommoditiesParameter{Interval2: p.Interval}
	inner, err := t.innerCommoditiesRequest(pp, "ALL_COMMODITIES")

	if err != nil {
		return nil, err
	}

	result := &AllCommodities{*inner}
	return result, nil
}
