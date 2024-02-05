package alphavantage

type FundamentalsClient struct {
	innerClient
}

// We offer the following set of fundamental data APIs in various temporal dimensions covering key financial metrics,
// income statements, balance sheets, cash flow, and other fundamental data points.
func NewFundamentalsClient(apikey string) *FundamentalsClient {
	return &FundamentalsClient{
		innerClient{
			httpClient: newHttpClient(),
			apikey:     apikey,
		},
	}
}

// parameter for OVERVIEW API
type CompanyOverviewParameter struct {
	// The symbol of the ticker of your choice. For example: symbol=IBM.
	Symbol string `json:"symbol"`
}

type CompanyOverview struct {
	Symbol                     string `json:"Symbol"`
	AssetType                  string `json:"AssetType"`
	Name                       string `json:"Name"`
	Description                string `json:"Description"`
	Cik                        string `json:"CIK"`
	Exchange                   string `json:"Exchange"`
	Currency                   string `json:"Currency"`
	Country                    string `json:"Country"`
	Sector                     string `json:"Sector"`
	Industry                   string `json:"Industry"`
	Address                    string `json:"Address"`
	FiscalYearEnd              string `json:"FiscalYearEnd"`
	LatestQuarter              string `json:"LatestQuarter"`
	MarketCapitalization       string `json:"MarketCapitalization"`
	Ebitda                     string `json:"EBITDA"`
	PERatio                    string `json:"PERatio"`
	PEGRatio                   string `json:"PEGRatio"`
	BookValue                  string `json:"BookValue"`
	DividendPerShare           string `json:"DividendPerShare"`
	DividendYield              string `json:"DividendYield"`
	Eps                        string `json:"EPS"`
	RevenuePerShareTTM         string `json:"RevenuePerShareTTM"`
	ProfitMargin               string `json:"ProfitMargin"`
	OperatingMarginTTM         string `json:"OperatingMarginTTM"`
	ReturnOnAssetsTTM          string `json:"ReturnOnAssetsTTM"`
	ReturnOnEquityTTM          string `json:"ReturnOnEquityTTM"`
	RevenueTTM                 string `json:"RevenueTTM"`
	GrossProfitTTM             string `json:"GrossProfitTTM"`
	DilutedEPSTTM              string `json:"DilutedEPSTTM"`
	QuarterlyEarningsGrowthYOY string `json:"QuarterlyEarningsGrowthYOY"`
	QuarterlyRevenueGrowthYOY  string `json:"QuarterlyRevenueGrowthYOY"`
	AnalystTargetPrice         string `json:"AnalystTargetPrice"`
	TrailingPE                 string `json:"TrailingPE"`
	ForwardPE                  string `json:"ForwardPE"`
	PriceToSalesRatioTTM       string `json:"PriceToSalesRatioTTM"`
	PriceToBookRatio           string `json:"PriceToBookRatio"`
	EVToRevenue                string `json:"EVToRevenue"`
	EVToEBITDA                 string `json:"EVToEBITDA"`
	Beta                       string `json:"Beta"`
	The52WeekHigh              string `json:"52WeekHigh"`
	The52WeekLow               string `json:"52WeekLow"`
	The50DayMovingAverage      string `json:"50DayMovingAverage"`
	The200DayMovingAverage     string `json:"200DayMovingAverage"`
	SharesOutstanding          string `json:"SharesOutstanding"`
	DividendDate               string `json:"DividendDate"`
	ExDividendDate             string `json:"ExDividendDate"`
}

func (t *FundamentalsClient) CompanyOverview(p CompanyOverviewParameter) (*CompanyOverview, error) {
	dic := make(map[string]string)
	dic["function"] = "OVERVIEW"
	dic["symbol"] = p.Symbol

	path := t.createQuerytUrl(dic)
	result := &CompanyOverview{}

	err := t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
