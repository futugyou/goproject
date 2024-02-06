package alphavantage

import "time"

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

// parameter for INCOME_STATEMENT API
type IncomeStatementParameter struct {
	// The symbol of the ticker of your choice. For example: symbol=IBM.
	Symbol string `json:"symbol"`
}

type IncomeStatement struct {
	Symbol           string                  `json:"symbol"`
	AnnualReports    []IncomeStatementReport `json:"annualReports"`
	QuarterlyReports []IncomeStatementReport `json:"quarterlyReports"`
}

type IncomeStatementReport struct {
	FiscalDateEnding                  string `json:"fiscalDateEnding"`
	ReportedCurrency                  string `json:"reportedCurrency"`
	GrossProfit                       string `json:"grossProfit"`
	TotalRevenue                      string `json:"totalRevenue"`
	CostOfRevenue                     string `json:"costOfRevenue"`
	CostofGoodsAndServicesSold        string `json:"costofGoodsAndServicesSold"`
	OperatingIncome                   string `json:"operatingIncome"`
	SellingGeneralAndAdministrative   string `json:"sellingGeneralAndAdministrative"`
	ResearchAndDevelopment            string `json:"researchAndDevelopment"`
	OperatingExpenses                 string `json:"operatingExpenses"`
	InvestmentIncomeNet               string `json:"investmentIncomeNet"`
	NetInterestIncome                 string `json:"netInterestIncome"`
	InterestIncome                    string `json:"interestIncome"`
	InterestExpense                   string `json:"interestExpense"`
	NonInterestIncome                 string `json:"nonInterestIncome"`
	OtherNonOperatingIncome           string `json:"otherNonOperatingIncome"`
	Depreciation                      string `json:"depreciation"`
	DepreciationAndAmortization       string `json:"depreciationAndAmortization"`
	IncomeBeforeTax                   string `json:"incomeBeforeTax"`
	IncomeTaxExpense                  string `json:"incomeTaxExpense"`
	InterestAndDebtExpense            string `json:"interestAndDebtExpense"`
	NetIncomeFromContinuingOperations string `json:"netIncomeFromContinuingOperations"`
	ComprehensiveIncomeNetOfTax       string `json:"comprehensiveIncomeNetOfTax"`
	Ebit                              string `json:"ebit"`
	Ebitda                            string `json:"ebitda"`
	NetIncome                         string `json:"netIncome"`
}

// This API returns the annual and quarterly income statements for the company of interest,
// with normalized fields mapped to GAAP and IFRS taxonomies of the SEC.
// Data is generally refreshed on the same day a company reports its latest earnings and financials.
func (t *FundamentalsClient) IncomeStatement(p IncomeStatementParameter) (*IncomeStatement, error) {
	dic := make(map[string]string)
	dic["function"] = "INCOME_STATEMENT"
	dic["symbol"] = p.Symbol

	path := t.createQuerytUrl(dic)
	result := &IncomeStatement{}

	err := t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// parameter for BALANCE_SHEET API
type BalanceSheetParameter struct {
	// The symbol of the ticker of your choice. For example: symbol=IBM.
	Symbol string `json:"symbol"`
}

type BalanceSheet struct {
	Symbol           string               `json:"symbol"`
	AnnualReports    []BalanceSheetReport `json:"annualReports"`
	QuarterlyReports []BalanceSheetReport `json:"quarterlyReports"`
}

type BalanceSheetReport struct {
	FiscalDateEnding                       string `json:"fiscalDateEnding"`
	ReportedCurrency                       string `json:"reportedCurrency"`
	TotalAssets                            string `json:"totalAssets"`
	TotalCurrentAssets                     string `json:"totalCurrentAssets"`
	CashAndCashEquivalentsAtCarryingValue  string `json:"cashAndCashEquivalentsAtCarryingValue"`
	CashAndShortTermInvestments            string `json:"cashAndShortTermInvestments"`
	Inventory                              string `json:"inventory"`
	CurrentNetReceivables                  string `json:"currentNetReceivables"`
	TotalNonCurrentAssets                  string `json:"totalNonCurrentAssets"`
	PropertyPlantEquipment                 string `json:"propertyPlantEquipment"`
	AccumulatedDepreciationAmortizationPPE string `json:"accumulatedDepreciationAmortizationPPE"`
	IntangibleAssets                       string `json:"intangibleAssets"`
	IntangibleAssetsExcludingGoodwill      string `json:"intangibleAssetsExcludingGoodwill"`
	Goodwill                               string `json:"goodwill"`
	Investments                            string `json:"investments"`
	LongTermInvestments                    string `json:"longTermInvestments"`
	ShortTermInvestments                   string `json:"shortTermInvestments"`
	OtherCurrentAssets                     string `json:"otherCurrentAssets"`
	OtherNonCurrentAssets                  string `json:"otherNonCurrentAssets"`
	TotalLiabilities                       string `json:"totalLiabilities"`
	TotalCurrentLiabilities                string `json:"totalCurrentLiabilities"`
	CurrentAccountsPayable                 string `json:"currentAccountsPayable"`
	DeferredRevenue                        string `json:"deferredRevenue"`
	CurrentDebt                            string `json:"currentDebt"`
	ShortTermDebt                          string `json:"shortTermDebt"`
	TotalNonCurrentLiabilities             string `json:"totalNonCurrentLiabilities"`
	CapitalLeaseObligations                string `json:"capitalLeaseObligations"`
	LongTermDebt                           string `json:"longTermDebt"`
	CurrentLongTermDebt                    string `json:"currentLongTermDebt"`
	LongTermDebtNoncurrent                 string `json:"longTermDebtNoncurrent"`
	ShortLongTermDebtTotal                 string `json:"shortLongTermDebtTotal"`
	OtherCurrentLiabilities                string `json:"otherCurrentLiabilities"`
	OtherNonCurrentLiabilities             string `json:"otherNonCurrentLiabilities"`
	TotalShareholderEquity                 string `json:"totalShareholderEquity"`
	TreasuryStock                          string `json:"treasuryStock"`
	RetainedEarnings                       string `json:"retainedEarnings"`
	CommonStock                            string `json:"commonStock"`
	CommonStockSharesOutstanding           string `json:"commonStockSharesOutstanding"`
}

// This API returns the annual and quarterly balance sheets for the company of interest,
// with normalized fields mapped to GAAP and IFRS taxonomies of the SEC.
// Data is generally refreshed on the same day a company reports its latest earnings and financials.
func (t *FundamentalsClient) BalanceSheet(p BalanceSheetParameter) (*BalanceSheet, error) {
	dic := make(map[string]string)
	dic["function"] = "BALANCE_SHEET"
	dic["symbol"] = p.Symbol

	path := t.createQuerytUrl(dic)
	result := &BalanceSheet{}

	err := t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// parameter for CASH_FLOW API
type CashFlowParameter struct {
	// The symbol of the ticker of your choice. For example: symbol=IBM.
	Symbol string `json:"symbol"`
}

type CashFlow struct {
	Symbol           string           `json:"symbol"`
	AnnualReports    []CashFlowReport `json:"annualReports"`
	QuarterlyReports []CashFlowReport `json:"quarterlyReports"`
}

type CashFlowReport struct {
	FiscalDateEnding                                          string `json:"fiscalDateEnding"`
	ReportedCurrency                                          string `json:"reportedCurrency"`
	OperatingCashflow                                         string `json:"operatingCashflow"`
	PaymentsForOperatingActivities                            string `json:"paymentsForOperatingActivities"`
	ProceedsFromOperatingActivities                           string `json:"proceedsFromOperatingActivities"`
	ChangeInOperatingLiabilities                              string `json:"changeInOperatingLiabilities"`
	ChangeInOperatingAssets                                   string `json:"changeInOperatingAssets"`
	DepreciationDepletionAndAmortization                      string `json:"depreciationDepletionAndAmortization"`
	CapitalExpenditures                                       string `json:"capitalExpenditures"`
	ChangeInReceivables                                       string `json:"changeInReceivables"`
	ChangeInInventory                                         string `json:"changeInInventory"`
	ProfitLoss                                                string `json:"profitLoss"`
	CashflowFromInvestment                                    string `json:"cashflowFromInvestment"`
	CashflowFromFinancing                                     string `json:"cashflowFromFinancing"`
	ProceedsFromRepaymentsOfShortTermDebt                     string `json:"proceedsFromRepaymentsOfShortTermDebt"`
	PaymentsForRepurchaseOfCommonStock                        string `json:"paymentsForRepurchaseOfCommonStock"`
	PaymentsForRepurchaseOfEquity                             string `json:"paymentsForRepurchaseOfEquity"`
	PaymentsForRepurchaseOfPreferredStock                     string `json:"paymentsForRepurchaseOfPreferredStock"`
	DividendPayout                                            string `json:"dividendPayout"`
	DividendPayoutCommonStock                                 string `json:"dividendPayoutCommonStock"`
	DividendPayoutPreferredStock                              string `json:"dividendPayoutPreferredStock"`
	ProceedsFromIssuanceOfCommonStock                         string `json:"proceedsFromIssuanceOfCommonStock"`
	ProceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet string `json:"proceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet"`
	ProceedsFromIssuanceOfPreferredStock                      string `json:"proceedsFromIssuanceOfPreferredStock"`
	ProceedsFromRepurchaseOfEquity                            string `json:"proceedsFromRepurchaseOfEquity"`
	ProceedsFromSaleOfTreasuryStock                           string `json:"proceedsFromSaleOfTreasuryStock"`
	ChangeInCashAndCashEquivalents                            string `json:"changeInCashAndCashEquivalents"`
	ChangeInExchangeRate                                      string `json:"changeInExchangeRate"`
	NetIncome                                                 string `json:"netIncome"`
}

// This API returns the annual and quarterly cash flow for the company of interest,
// with normalized fields mapped to GAAP and IFRS taxonomies of the SEC.
// Data is generally refreshed on the same day a company reports its latest earnings and financials.
func (t *FundamentalsClient) CashFlow(p CashFlowParameter) (*CashFlow, error) {
	dic := make(map[string]string)
	dic["function"] = "CASH_FLOW"
	dic["symbol"] = p.Symbol

	path := t.createQuerytUrl(dic)
	result := &CashFlow{}

	err := t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// parameter for EARNINGS API
type EarningsParameter struct {
	// The symbol of the ticker of your choice. For example: symbol=IBM.
	Symbol string `json:"symbol"`
}

type Earnings struct {
	Symbol            string             `json:"symbol"`
	AnnualEarnings    []AnnualEarning    `json:"annualEarnings"`
	QuarterlyEarnings []QuarterlyEarning `json:"quarterlyEarnings"`
}

type AnnualEarning struct {
	FiscalDateEnding string `json:"fiscalDateEnding"`
	ReportedEPS      string `json:"reportedEPS"`
}

type QuarterlyEarning struct {
	FiscalDateEnding   string `json:"fiscalDateEnding"`
	ReportedDate       string `json:"reportedDate"`
	ReportedEPS        string `json:"reportedEPS"`
	EstimatedEPS       string `json:"estimatedEPS"`
	Surprise           string `json:"surprise"`
	SurprisePercentage string `json:"surprisePercentage"`
}

// This API returns the annual and quarterly earnings (EPS) for the company of interest.
// Quarterly data also includes analyst estimates and surprise metrics.
func (t *FundamentalsClient) Earnings(p EarningsParameter) (*Earnings, error) {
	dic := make(map[string]string)
	dic["function"] = "EARNINGS"
	dic["symbol"] = p.Symbol

	path := t.createQuerytUrl(dic)
	result := &Earnings{}

	err := t.httpClient.getJson(path, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// parameter for LISTING_STATUS API
type ListingStatusParameter struct {
	// If no date is set, the API endpoint will return a list of active or delisted symbols as of the latest trading day.
	// If a date is set, the API endpoint will "travel back" in time and return a list of active or delisted symbols on that particular date in history.
	// Any YYYY-MM-DD date later than 2010-01-01 is supported. For example, date=2013-08-03
	Date string `json:"date"`
	// By default, state=active and the API will return a list of actively traded stocks and ETFs.
	// Set state=delisted to query a list of delisted assets.
	State string `json:"state"`
}

// symbol,name,exchange,assetType,ipoDate,delistingDate,status
type ListingStatus struct {
	Symbol        string    `json:"symbol" csv:"symbol"`
	Name          string    `json:"name" csv:"name"`
	Exchange      string    `json:"exchange" csv:"exchange"`
	AssetType     string    `json:"assetType" csv:"assetType"`
	IpoDate       time.Time `json:"ipoDate" csv:"ipoDate"`
	DelistingDate time.Time `json:"delistingDate" csv:"delistingDate"`
	Status        string    `json:"status" csv:"status"`
}

// This API returns the annual and quarterly earnings (EPS) for the company of interest.
// Quarterly data also includes analyst estimates and surprise metrics.
func (t *FundamentalsClient) ListingStatus(p ListingStatusParameter) ([]ListingStatus, error) {
	dic := make(map[string]string)
	dic["function"] = "LISTING_STATUS"
	dic["date"] = p.Date
	dic["state"] = p.State

	path := t.createQuerytUrl(dic)
	result := make([]ListingStatus, 0)

	err := t.httpClient.getCsvByUtil(path, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// parameter for EARNINGS_CALENDAR API
type EarningsCalendarParameter struct {
	// By default, no symbol will be set for this API. When no symbol is set, the API endpoint will return the full list of company earnings scheduled.
	// If a symbol is set, the API endpoint will return the expected earnings for that specific symbol. For example, symbol=IBM
	Symbol string `json:"symbol"`
	// By default, horizon=3month and the API will return a list of expected company earnings in the next 3 months.
	// You may set horizon=6month or horizon=12month to query the earnings scheduled for the next 6 months or 12 months, respectively.
	Horizon string `json:"horizon"`
}

// symbol,name,reportDate,fiscalDateEnding,estimate,currency
type EarningsCalendar struct {
	Symbol           string    `json:"symbol" csv:"symbol"`
	Name             string    `json:"name" csv:"name"`
	ReportDate       time.Time `json:"reportDate" csv:"reportDate"`
	FiscalDateEnding time.Time `json:"fiscalDateEnding" csv:"fiscalDateEnding"`
	Estimate         float64   `json:"estimate" csv:"estimate"`
	Currency         string    `json:"currency" csv:"currency"`
}

// This API returns the annual and quarterly earnings (EPS) for the company of interest.
// Quarterly data also includes analyst estimates and surprise metrics.
func (t *FundamentalsClient) EarningsCalendar(p EarningsCalendarParameter) ([]EarningsCalendar, error) {
	dic := make(map[string]string)
	dic["function"] = "EARNINGS_CALENDAR"
	dic["symbol"] = p.Symbol
	dic["horizon"] = p.Horizon

	path := t.createQuerytUrl(dic)
	result := make([]EarningsCalendar, 0)

	err := t.httpClient.getCsvByUtil(path, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
