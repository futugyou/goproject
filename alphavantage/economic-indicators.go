package alphavantage

import "github.com/futugyou/alphavantage/enums"

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

type EconomicIndicatorsParameter struct {
	// By default, interval=monthly. Strings daily, weekly, and monthly are accepted.
	Interval enums.LongInterval           `json:"interval"`
	Function enums.EconomicIndicatorsType `json:"function"`
	Maturity enums.Maturity               `json:"maturity"`
}

func (p EconomicIndicatorsParameter) Validation() (map[string]string, error) {
	dic := make(map[string]string)
	dic["function"] = p.Function.String()
	if p.Interval != nil {
		dic["interval"] = p.Interval.String()
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

func (t *EconomicIndicatorsClient) GetEconomicIndicatorsDara(p EconomicIndicatorsParameter) (*EconomicIndicatorsResult, error) {
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
