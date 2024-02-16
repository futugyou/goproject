package stock

type StockEntity struct {
	Id                         string  `bson:"_id,omitempty"`
	Symbol                     string  `bson:"symbol"`
	Name                       string  `bson:"name"`
	Type                       string  `bson:"type"`
	Region                     string  `bson:"region"`
	MarketOpen                 string  `bson:"marketOpen"`
	MarketClose                string  `bson:"marketClose"`
	Timezone                   string  `bson:"timezone"`
	Currency                   string  `bson:"currency"`
	MatchScore                 float64 `bson:"matchScore"`
	AssetType                  string  `bson:"AssetType"`
	Description                string  `bson:"Description"`
	Cik                        string  `bson:"CIK"`
	Exchange                   string  `bson:"Exchange"`
	Country                    string  `bson:"Country"`
	Sector                     string  `bson:"Sector"`
	Industry                   string  `bson:"Industry"`
	Address                    string  `bson:"Address"`
	FiscalYearEnd              string  `bson:"FiscalYearEnd"`
	LatestQuarter              string  `bson:"LatestQuarter"`
	MarketCapitalization       string  `bson:"MarketCapitalization"`
	Ebitda                     string  `bson:"EBITDA"`
	PERatio                    string  `bson:"PERatio"`
	PEGRatio                   string  `bson:"PEGRatio"`
	BookValue                  string  `bson:"BookValue"`
	DividendPerShare           string  `bson:"DividendPerShare"`
	DividendYield              string  `bson:"DividendYield"`
	Eps                        string  `bson:"EPS"`
	RevenuePerShareTTM         string  `bson:"RevenuePerShareTTM"`
	ProfitMargin               string  `bson:"ProfitMargin"`
	OperatingMarginTTM         string  `bson:"OperatingMarginTTM"`
	ReturnOnAssetsTTM          string  `bson:"ReturnOnAssetsTTM"`
	ReturnOnEquityTTM          string  `bson:"ReturnOnEquityTTM"`
	RevenueTTM                 string  `bson:"RevenueTTM"`
	GrossProfitTTM             string  `bson:"GrossProfitTTM"`
	DilutedEPSTTM              string  `bson:"DilutedEPSTTM"`
	QuarterlyEarningsGrowthYOY string  `bson:"QuarterlyEarningsGrowthYOY"`
	QuarterlyRevenueGrowthYOY  string  `bson:"QuarterlyRevenueGrowthYOY"`
	AnalystTargetPrice         string  `bson:"AnalystTargetPrice"`
	TrailingPE                 string  `bson:"TrailingPE"`
	ForwardPE                  string  `bson:"ForwardPE"`
	PriceToSalesRatioTTM       string  `bson:"PriceToSalesRatioTTM"`
	PriceToBookRatio           string  `bson:"PriceToBookRatio"`
	EVToRevenue                string  `bson:"EVToRevenue"`
	EVToEBITDA                 string  `bson:"EVToEBITDA"`
	Beta                       string  `bson:"Beta"`
	The52WeekHigh              string  `bson:"52WeekHigh"`
	The52WeekLow               string  `bson:"52WeekLow"`
	The50DayMovingAverage      string  `bson:"50DayMovingAverage"`
	The200DayMovingAverage     string  `bson:"200DayMovingAverage"`
	SharesOutstanding          string  `bson:"SharesOutstanding"`
	DividendDate               string  `bson:"DividendDate"`
	ExDividendDate             string  `bson:"ExDividendDate"`
}

func (StockEntity) GetType() string {
	return "stocks"
}
