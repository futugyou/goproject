package income

type IncomeEntity struct {
	Id                                string `bson:"_id"`
	Symbol                            string `bson:"symbol"`
	DataType                          string `bson:"dataType"`
	FiscalDateEnding                  string `bson:"fiscalDateEnding"`
	ReportedCurrency                  string `bson:"reportedCurrency"`
	GrossProfit                       string `bson:"grossProfit"`
	TotalRevenue                      string `bson:"totalRevenue"`
	CostOfRevenue                     string `bson:"costOfRevenue"`
	CostofGoodsAndServicesSold        string `bson:"costofGoodsAndServicesSold"`
	OperatingIncome                   string `bson:"operatingIncome"`
	SellingGeneralAndAdministrative   string `bson:"sellingGeneralAndAdministrative"`
	ResearchAndDevelopment            string `bson:"researchAndDevelopment"`
	OperatingExpenses                 string `bson:"operatingExpenses"`
	InvestmentIncomeNet               string `bson:"investmentIncomeNet"`
	NetInterestIncome                 string `bson:"netInterestIncome"`
	InterestIncome                    string `bson:"interestIncome"`
	InterestExpense                   string `bson:"interestExpense"`
	NonInterestIncome                 string `bson:"nonInterestIncome"`
	OtherNonOperatingIncome           string `bson:"otherNonOperatingIncome"`
	Depreciation                      string `bson:"depreciation"`
	DepreciationAndAmortization       string `bson:"depreciationAndAmortization"`
	IncomeBeforeTax                   string `bson:"incomeBeforeTax"`
	IncomeTaxExpense                  string `bson:"incomeTaxExpense"`
	InterestAndDebtExpense            string `bson:"interestAndDebtExpense"`
	NetIncomeFromContinuingOperations string `bson:"netIncomeFromContinuingOperations"`
	ComprehensiveIncomeNetOfTax       string `bson:"comprehensiveIncomeNetOfTax"`
	Ebit                              string `bson:"ebit"`
	Ebitda                            string `bson:"ebitda"`
	NetIncome                         string `bson:"netIncome"`
}

func (IncomeEntity) GetType() string {
	return "Incomes"
}
