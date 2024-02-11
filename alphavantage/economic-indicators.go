package alphavantage

import (
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
	Name     string  `json:"name"`
	Interval string  `json:"interval"`
	Unit     string  `json:"unit"`
	Data     []Datum `json:"data"`
}

func (t *EconomicIndicatorsClient) GetEconomicIndicatorsDara(p innerEconomicIndicatorsParameter) (*EconomicIndicatorsResult, error) {
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

	return result, nil
}

type TreasuryYieldParameter struct {
	Interval enums.EconomicTreasuryInterval `json:"interval"`
	Maturity enums.Maturity                 `json:"maturity"`
}

// This API returns the annual and quarterly Real GDP of the United States.
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

	return result, nil
}

type InterestRateParameter struct {
	Interval enums.EconomicFundsInterval `json:"interval"` 
}

// This API returns the annual and quarterly Real GDP of the United States.
func (t *EconomicIndicatorsClient) InterestRate(p InterestRateParameter) (*EconomicIndicatorsResult, error) {
	pp := innerEconomicIndicatorsParameter{FundsInterval: p.Interval,  Function: functions.FederalFunds}
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

	return result, nil
}
