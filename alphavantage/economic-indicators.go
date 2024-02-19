package alphavantage

import (
	"fmt"

	"github.com/futugyou/alphavantage/enums"
	"github.com/futugyou/alphavantage/functions"
)

type EconomicIndicatorsClient struct {
	innerClient
}

// APIs under this section provide key US economic indicators frequently used for investment strategy formulation and application development.
func NewEconomicIndicatorsClient(apikey string) *EconomicIndicatorsClient {
	return &EconomicIndicatorsClient{
		innerClient{
			httpClient: newHttpClient(),
			apikey:     apikey,
		},
	}
}

type innerEconomicIndicatorsParameter struct {
	Function         functions.EconomicIndicatorsType
	Maturity         enums.Maturity
	GdpInterval      enums.EconomicGdpInterval
	TreasuryInterval enums.EconomicTreasuryInterval
	FundsInterval    enums.EconomicFundsInterval
	CPIInterval      enums.EconomicCPIInterval
}

func (p innerEconomicIndicatorsParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = p.Function.String()
	if p.GdpInterval != nil {
		dic["interval"] = p.GdpInterval.String()
	}
	if p.TreasuryInterval != nil {
		dic["interval"] = p.TreasuryInterval.String()
	}
	if p.FundsInterval != nil {
		dic["interval"] = p.FundsInterval.String()
	}
	if p.CPIInterval != nil {
		dic["interval"] = p.CPIInterval.String()
	}
	if p.Maturity != nil {
		dic["maturity"] = p.Maturity.String()
	}

	dic["datatype"] = "json"
	return dic, nil
}

type EconomicIndicatorsResult struct {
	Name        string  `json:"name"`
	Interval    string  `json:"interval"`
	Unit        string  `json:"unit"`
	Data        []Datum `json:"data"`
	Information string  `json:"Information"`
}

func (t *EconomicIndicatorsClient) GetEconomicIndicatorsData(p innerEconomicIndicatorsParameter) (*EconomicIndicatorsResult, error) {
	dic, err := p.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

type RealGdpParameter struct {
	Interval enums.EconomicGdpInterval `json:"interval"`
}

// This API returns the annual and quarterly Real GDP of the United States.
func (t *EconomicIndicatorsClient) RealGdp(p RealGdpParameter) (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{GdpInterval: p.Interval, Function: functions.RealGDP}
	dic, err := pp.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

// This API returns the quarterly Real GDP per Capita data of the United States.
func (t *EconomicIndicatorsClient) RealGdpPerCapita() (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{Function: functions.RealGDPPerCapita}
	dic, err := pp.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

type TreasuryYieldParameter struct {
	Interval enums.EconomicTreasuryInterval `json:"interval"`
	Maturity enums.Maturity                 `json:"maturity"`
}

// This API returns the daily, weekly, and monthly US treasury yield of a given maturity timeline (e.g., 5 year, 30 year, etc).
func (t *EconomicIndicatorsClient) TreasuryYield(p TreasuryYieldParameter) (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{TreasuryInterval: p.Interval, Maturity: p.Maturity, Function: functions.TreasuryYield}
	dic, err := pp.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

type InterestRateParameter struct {
	Interval enums.EconomicFundsInterval `json:"interval"`
}

// This API returns the daily, weekly, and monthly federal funds rate (interest rate) of the United States.
func (t *EconomicIndicatorsClient) InterestRate(p InterestRateParameter) (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{FundsInterval: p.Interval, Function: functions.FederalFunds}
	dic, err := pp.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

type CPIParameter struct {
	Interval enums.EconomicCPIInterval `json:"interval"`
}

// This API returns the monthly and semiannual consumer price index (CPI) of the United States.
// CPI is widely regarded as the barometer of inflation levels in the broader economy.
func (t *EconomicIndicatorsClient) CPI(p CPIParameter) (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{CPIInterval: p.Interval, Function: functions.CPI}
	dic, err := pp.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

// This API returns the annual inflation rates (consumer prices) of the United States.
func (t *EconomicIndicatorsClient) Inflation() (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{Function: functions.Inflation}
	dic, err := pp.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

// This API returns the monthly Advance Retail Sales: Retail Trade data of the United States.
func (t *EconomicIndicatorsClient) RetailSales() (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{Function: functions.RetailSales}
	dic, err := pp.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

// This API returns the monthly manufacturers' new orders of durable goods in the United States.
func (t *EconomicIndicatorsClient) DurableGoods() (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{Function: functions.DurableGoods}
	dic, err := pp.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

// This API returns the monthly unemployment data of the United States.
// The unemployment rate represents the number of unemployed as a percentage of the labor force.
// Labor force data are restricted to people 16 years of age and older, who currently reside in 1 of the 50 states or the District of Columbia,
// who do not reside in institutions (e.g., penal and mental facilities, homes for the aged), and who are not on active duty in the Armed Forces (source).
func (t *EconomicIndicatorsClient) UnemploymentRate() (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{Function: functions.UnemploymentRate}
	dic, err := pp.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}

// This API returns the monthly US All Employees:
// Total Nonfarm (commonly known as Total Nonfarm Payroll),
// a measure of the number of U.S. workers in the economy that excludes proprietors, private household employees, unpaid volunteers, farm employees,
// and the unincorporated self-employed.
func (t *EconomicIndicatorsClient) NonfarmPayroll() (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{Function: functions.NonfarmPayroll}
	dic, err := pp.Validation()
	if err != nil {
		return nil, err
	}

	path := t.createQuerytUrl(dic)
	result := &EconomicIndicatorsResult{}
	err = t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	if len(result.Information) > 0 {
		return nil, fmt.Errorf(result.Information)
	}

	return result, nil
}
